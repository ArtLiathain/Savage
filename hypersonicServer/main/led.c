#include "../include/led.h"
#include <driver/ledc.h>
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "../include/shared_memory.h"
#include "esp_log.h"
#include "../include/hypersonic.h"

#define TIMER LEDC_TIMER_0
#define MODE LEDC_LOW_SPEED_MODE

typedef enum
{
    STATE_MIDDLE, // 0
    STATE_CLOSE,  // 1
    STATE_FAR,    // 2
    STATE_MAX     // 3
} DistanceState;

#define CLOSE_THRESHOLD 15
#define FAR_THRESHOLD 30

DistanceState determine_state(int distance)
{
    if (distance < CLOSE_THRESHOLD && distance > 0)
    {
        return STATE_CLOSE;
    }
    else if (distance > FAR_THRESHOLD)
    {
        return STATE_FAR;
    }
    else
    {
        return STATE_MIDDLE;
    }
}

static DistanceState current_state = STATE_MIDDLE; // Start with an initial state
void handle_distance_state(uint8_t distance, shared_memory_t *shared_memory)
{

    DistanceState next_state = determine_state(distance);

    // Handle state transitions
    if (next_state != current_state)
    {
        printf("State changed from %d to %d\n", current_state, next_state);

        // Optional: Perform actions when entering a new state
        switch (next_state)
        {
        case STATE_CLOSE:
            update_led_duty(shared_memory->leds_config[0], 1 << CONFIG_LED_SHADE_BITS);
            update_led_duty(shared_memory->leds_config[1], 0);
            ESP_LOGI("LED", "Object is at a close distance");

            break;
        case STATE_FAR:
            update_led_duty(shared_memory->leds_config[1], 1 << CONFIG_LED_SHADE_BITS);
            update_led_duty(shared_memory->leds_config[0], 0);
            ESP_LOGI("LED", "Object is at a far distance");

            break;
        case STATE_MIDDLE:
            update_led_duty(shared_memory->leds_config[1], 0);
            update_led_duty(shared_memory->leds_config[0], 0);
            ESP_LOGI("LED", "Object is at a middle distance");
            break;
        default:
            break;
        }

        // Update the current state
        current_state = next_state;
    }
}

void init_led_timer(int frequency, int bits)
{
    ledc_timer_config_t ledc_timer = {
        .speed_mode = MODE,
        .duty_resolution = bits,
        .timer_num = TIMER,
        .freq_hz = frequency,
        .clk_cfg = LEDC_AUTO_CLK};
    ESP_ERROR_CHECK(ledc_timer_config(&ledc_timer));
}

void init_led_on_pin(led_config_t config)
{

    ledc_channel_config_t ledc_channel1 = {
        .speed_mode = MODE,
        .channel = config.channel,
        .timer_sel = TIMER,
        .intr_type = LEDC_INTR_DISABLE,
        .gpio_num = config.gpio_pin,
        .duty = 0, // Set duty to 0%
        .hpoint = 0};
    ESP_ERROR_CHECK(ledc_channel_config(&ledc_channel1));
}

void update_led_duty(led_config_t config, int shade)
{
    ledc_set_duty(MODE, config.channel, shade);
    ledc_update_duty(MODE, config.channel);
}

void led_state_machine_task(void *pvParameters)
{
    shared_memory_t *shared_memory = (shared_memory_t *)pvParameters;
    uint32_t notificationValue;
    TaskHandle_t led_breathing_handle = NULL; // Declare a TaskHandle_t variable to hold the task handle

    for (;;)
    {
        notificationValue = ulTaskNotifyTake(pdTRUE, portMAX_DELAY);

        if (notificationValue == 1)
        {
            uint8_t distance = get_distance(shared_memory->hyper_config);
            handle_distance_state(distance, shared_memory);
            circular_buffer_push(shared_memory->led_buffer, current_state);
        }

        if (notificationValue == 2)
        {
            if (led_breathing_handle == NULL)
            {
                xTaskCreate(led_breathe_task, "led_breathe_task", 4096, shared_memory, 1, &led_breathing_handle);
            }
        }
        if (notificationValue == 3)
        {
            vTaskDelete(led_breathing_handle);
            led_breathing_handle = NULL;
        }
    }
}

void led_breathe_task(void *pvParameters)
{
    shared_memory_t *shared_memory = (shared_memory_t *)pvParameters;

    int shade = 0;
    int increasing = 1;
    for (;;)

    {
        for (int index = 0; index < 2; index++)
        {
            update_led_duty(shared_memory->leds_config[index], shade);
        }

        if (increasing == 1)
        {
            shade++;
            if (shade >= (1 << CONFIG_LED_SHADE_BITS) - 1)
            {
                increasing = 0;
            }
        }
        else
        {
            shade--;
            if (shade <= 0)
            {
                increasing = 1;
            }
        }

        vTaskDelay(CONFIG_LED_BLINK_SPEED / portTICK_PERIOD_MS);
    }
}