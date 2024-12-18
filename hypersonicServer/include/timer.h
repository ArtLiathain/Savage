#ifndef TIMER
#define TIMER
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"

void init_timer(int group, int timer, bool auto_reload, int timer_interval_sec, TaskHandle_t *handle);

#endif
