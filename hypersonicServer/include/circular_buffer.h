#ifndef CIRCULAR_BUFFER_H
#define CIRCULAR_BUFFER_H

#include <stdint.h> 
#include <stddef.h>
#include <stdbool.h>
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/semphr.h"


typedef struct
{
    float *buffer;
    size_t head;
    size_t tail;
    size_t max_size;
    bool empty;
    SemaphoreHandle_t*  mutex;
    
} CircularBuffer;


size_t circular_buffer_to_byte_array(CircularBuffer *cb, uint8_t amount ,uint8_t *byte_array);

bool circular_buffer_push(CircularBuffer *cb, uint8_t value);

uint8_t circular_buffer_pop(CircularBuffer *cb);

uint8_t circular_buffer_read_top(CircularBuffer *cb);

void circular_buffer_print(CircularBuffer *cb);

#endif
