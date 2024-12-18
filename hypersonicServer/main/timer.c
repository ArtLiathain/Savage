#include <stdio.h>
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "driver/timer.h"
#include "esp_system.h"
#include "esp_log.h"
#include "freertos/semphr.h"

#define TIMER_BASE_CLK 80000000
#define TIMER_DIVIDER 16                             // Timer clock divider
#define TIMER_SCALE (TIMER_BASE_CLK / TIMER_DIVIDER) // Convert to seconds

static bool IRAM_ATTR timer_group_isr_callback(void *args)
{
    BaseType_t high_task_awoken = pdFALSE;
    TaskHandle_t *taskHandle = (TaskHandle_t *)args;

    if (taskHandle != NULL && *taskHandle != NULL)
    {
        xTaskNotifyFromISR(*taskHandle, 1, eSetBits, &high_task_awoken);
    }

    return high_task_awoken == pdTRUE; // return whether we need to yield at the end of ISR
}

void init_timer(int group, int timer, bool auto_reload, int timer_interval_sec, TaskHandle_t *handle)
{
    /* Select and initialize basic parameters of the timer */
    timer_config_t config = {
        .divider = TIMER_DIVIDER,
        .counter_dir = TIMER_COUNT_UP,
        .counter_en = TIMER_PAUSE,
        .alarm_en = TIMER_ALARM_EN,
        .auto_reload = auto_reload,
    }; // default clock source is APB
    timer_init(group, timer, &config);

    /* Timer's counter will initially start from value below.
       Also, if auto_reload is set, this value will be automatically reload on alarm */
    timer_set_counter_value(group, timer, 0);

    /* Configure the alarm value and the interrupt on alarm. */
    timer_set_alarm_value(group, timer, timer_interval_sec * TIMER_SCALE);
    timer_enable_intr(group, timer);

    if (handle != NULL && *handle != NULL)
    {
        timer_isr_callback_add(group, timer, timer_group_isr_callback, handle, 0);
    }
    else
    {
        printf("Invalid task handle passed to init_timer\n");
    }

    timer_start(group, timer);
}
