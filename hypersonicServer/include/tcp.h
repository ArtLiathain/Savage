#ifndef TCP
#define TCP
#include "../include/shared_memory.h"
void tcp_server_task(void *pvParameters);
void wifi_init_sta();

typedef struct
{
    shared_memory_t *shared_memory;
    TaskHandle_t *led_task_handle;

} tcp_shared_data_t;
#endif