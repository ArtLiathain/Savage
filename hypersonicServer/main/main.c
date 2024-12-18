#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "esp_adc/adc_oneshot.h"
#include <driver/ledc.h>
#include "driver/timer.h"
#include "esp_log.h"
#include "../include/hypersonic.h"
#include "../include/led.h"
#include "../include/tcp.h"
#include "../include/circular_buffer.h"
#include "../include/timer.h"
#include "../include/shared_memory.h"
#include "esp_sntp.h" // New API for SNTP
#include "time.h"
#include "esp_netif.h"
#include "esp_wifi.h"  // Include for Wi-Fi functions
#include "nvs_flash.h" // Include for NVS functions
#include "freertos/semphr.h"
#include "sdkconfig.h"

#define CHANNEL1 LEDC_CHANNEL_0
#define CHANNEL2 LEDC_CHANNEL_1
#define FREQ 5000 // 5 kHz refresh
#define STACK_SIZE 4096

TaskHandle_t hypersonic_task_handle = NULL;
TaskHandle_t led_task_handle = NULL;
float hyper_buffer[CONFIG_BUFFER_SIZE] = {0};
float led_buffer[CONFIG_BUFFER_SIZE] = {0};
SemaphoreHandle_t hypersonic_mutex;
SemaphoreHandle_t led_mutex;
CircularBuffer hypersonic_buffer = {.buffer = hyper_buffer,
                                    .empty = true,
                                    .max_size = CONFIG_BUFFER_SIZE,
                                    .tail = 0,
                                    .head = 0,
                                    .mutex = &hypersonic_mutex};
CircularBuffer led_state_buffer = {.buffer = led_buffer,
                                   .empty = true,
                                   .max_size = CONFIG_BUFFER_SIZE,
                                   .tail = 0,
                                   .head = 0,
                                   .mutex = &led_mutex};

hypersonic_config hyper_config = {.channel = ADC_CHANNEL_6, .attenuation = ADC_ATTEN_DB_0};
shared_memory_t shared_mem;
tcp_shared_data_t shared_tcp_data;
led_config_t led1_config = {.channel = CHANNEL1, .gpio_pin = CONFIG_LED_PIN_YELLOW};
led_config_t led2_config = {.channel = CHANNEL2, .gpio_pin = CONFIG_LED_PIN_GREEN};
led_config_t leds[2];

void app_main()
{
    hypersonic_mutex = xSemaphoreCreateMutex();
    led_mutex = xSemaphoreCreateMutex();
    initialise_hypersonic(&hyper_config);
    init_led_timer(FREQ, CONFIG_LED_SHADE_BITS);
    init_led_on_pin(led1_config);
    init_led_on_pin(led2_config);
    esp_rom_gpio_pad_select_gpio(CONFIG_LED_PIN_GREEN);
    esp_rom_gpio_pad_select_gpio(CONFIG_LED_PIN_YELLOW);
    leds[0] = led1_config;
    leds[1] = led2_config;

    shared_mem.hyper_config = &hyper_config;
    shared_mem.sampling_frequency = CONFIG_TIMER_LENGTH;
    shared_mem.hypersonic_buffer = &hypersonic_buffer;
    shared_mem.led_buffer = &led_state_buffer;
    shared_mem.leds_config = &leds;

    xTaskCreate(sample_hypersonic_task, "sampling_hypersonic", STACK_SIZE, &shared_mem, 1, &hypersonic_task_handle);
    xTaskCreate(led_state_machine_task, "LED_state_management", STACK_SIZE, &shared_mem, 1, &led_task_handle);
    init_timer(TIMER_GROUP_0, TIMER_0, true, 2, &hypersonic_task_handle);
    vTaskDelay(1000 / portTICK_PERIOD_MS);
    init_timer(TIMER_GROUP_1, TIMER_0, true, 2, &led_task_handle);
    ESP_ERROR_CHECK(nvs_flash_init()); // Initialize NVS
    wifi_init_sta();                   // Initialize Wi-Fi
    shared_tcp_data.led_task_handle = &led_task_handle;
    shared_tcp_data.shared_memory = &shared_mem;
    xTaskCreate(&tcp_server_task, "tcp_server_task", STACK_SIZE, &shared_tcp_data, 5, NULL);
}
