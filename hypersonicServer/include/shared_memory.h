#ifndef SHARED_MEMORY
#define SHARED_MEMORY

#include "led.h"
#include <time.h>
#include "circular_buffer.h"
#include <stdint.h> 
#include "hypersonic.h"



typedef struct {
    CircularBuffer* hypersonic_buffer;
    CircularBuffer* led_buffer;
    uint8_t sampling_frequency;
    hypersonic_config* hyper_config;
    led_config_t*  leds_config;
} shared_memory_t;


#endif