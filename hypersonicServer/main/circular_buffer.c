#include "../include/circular_buffer.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include <string.h>

bool circular_buffer_push(CircularBuffer *cb, uint8_t value)
{
    if (xSemaphoreTake(*(cb->mutex), portMAX_DELAY))
    { // Lock the buffer

        cb->buffer[cb->head] = value;
        cb->head = (cb->head + 1) % cb->max_size;
        if (cb->tail == cb->head)
        {
            cb->tail = (cb->head + 1) % cb->max_size;
        }
        cb->empty = false;
        xSemaphoreGive(*(cb->mutex)); // Unlock the buffer

        return true;
    }
    return false;
}

uint8_t circular_buffer_pop(CircularBuffer *cb)
{
    if (xSemaphoreTake(*(cb->mutex), portMAX_DELAY))
    { // Lock the buffer
        if (cb->empty)
        {
            xSemaphoreGive(*(cb->mutex)); // Unlock the buffer
            return 0;
        }

        uint8_t val = cb->buffer[cb->tail];
        cb->tail = (cb->tail + 1) % cb->max_size;
        cb->empty = (cb->tail == cb->head);
        xSemaphoreGive(*(cb->mutex)); // Unlock the buffer
        return val;
    }

    return 0;
}

uint8_t circular_buffer_read_top(CircularBuffer *cb)
{
    if (xSemaphoreTake(*(cb->mutex), portMAX_DELAY))
    { // Lock the buffer
        if (cb->empty)
        {
            xSemaphoreGive(*(cb->mutex)); // Unlock the buffer
            return 0;
        }

        uint8_t val = cb->buffer[cb->head - 1];
        xSemaphoreGive(*(cb->mutex)); // Unlock the buffer
        return val;
    }
    return 0;
}

void circular_buffer_print(CircularBuffer *cb)
{
    if (cb->empty)
    {
        return;
    }
    int index = cb->tail;
    while (index != cb->head)
    {
        printf("Value: %f\n", cb->buffer[index]);
        index = (index + 1) % cb->max_size;
    }
}

size_t circular_buffer_to_byte_array(CircularBuffer *cb, uint8_t amount, uint8_t *byte_array)
{
    uint8_t count = 0;
    size_t byte_count = sizeof(uint8_t);

    if (cb->empty)
    {
        uint8_t empty_flag = 0x00;
        memcpy(byte_array, &empty_flag, sizeof(uint8_t));
        return byte_count; // No data to convert
    }

    while (count < amount && !cb->empty)
    {
        uint8_t val = circular_buffer_pop(cb);
        memcpy(byte_array + byte_count, &val, sizeof(uint8_t));
        byte_count += sizeof(uint8_t);
        count++;
    }
    memcpy(byte_array, &count, sizeof(uint8_t));
    return byte_count; // Return the number of bytes written to the byte array
}
