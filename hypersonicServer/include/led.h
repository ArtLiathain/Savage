#ifndef LED
#define LED
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"

typedef struct
{
    int channel;
    int gpio_pin;
} led_config_t;

void led_breathe_task(void *pvParam);
void init_led_timer(int frequency, int bits);
void init_led_on_pin(led_config_t config);
void update_led_duty(led_config_t config, int shade);
void led_state_machine_task(void *pvParameters);

#endif