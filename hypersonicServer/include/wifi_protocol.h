#ifndef WIFI_PROTOCOL
#define WIFI_PROTOCOL
#include "shared_memory.h"

void custom_receiver(int client_sock, shared_memory_t* shared_memory, TaskHandle_t* led_task_handle);

#endif