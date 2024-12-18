#include "../include/hypersonic.h"
#include "../include/shared_memory.h"
#include "freertos/FreeRTOS.h"
#include "esp_log.h"
#include "esp_adc/adc_oneshot.h"
#include <esp_err.h>

#define ADC_30CM 830 // ADC value for 30 cm

static adc_oneshot_unit_handle_t adc_handle;

void initialise_hypersonic(hypersonic_config *config)
{
    // Configure ADC unit
    adc_oneshot_unit_init_cfg_t init_cfg = {
        .unit_id = ADC_UNIT_1,
    };
    esp_err_t err = adc_oneshot_new_unit(&init_cfg, &adc_handle);
    if (err != ESP_OK)
    {
        printf("ADC unit initialization failed: %s\n", esp_err_to_name(err));
        return;
    }

    // Configure ADC channel
    adc_oneshot_chan_cfg_t chan_cfg = {
        .atten = config->attenuation, // ADC_ATTEN_DB_0, ADC_ATTEN_DB_2_5, etc.
        .bitwidth = ADC_BITWIDTH_DEFAULT,
    };
    err = adc_oneshot_config_channel(adc_handle, config->channel, &chan_cfg);
    if (err != ESP_OK)
    {
        printf("ADC channel configuration failed: %s\n", esp_err_to_name(err));
    }
}

int calculate_distance(int adc_value)
{
    int y;
    // Apply the formula y - 30 = 2(x - 830) / 77
    y = 30 + (2 * (adc_value - ADC_30CM)) / 77.0;

    return y;
}

int get_distance(hypersonic_config *config)
{
    int raw_value = 0;
    esp_err_t err = adc_oneshot_read(adc_handle, config->channel, &raw_value);
    if (err != ESP_OK)
    {
        printf("ADC read failed: %s\n", esp_err_to_name(err));
        return -1.0; // Indicate failure
    }
    return calculate_distance(raw_value);
}

void sample_hypersonic_task(void *params)
{
    shared_memory_t *shared_memory = (shared_memory_t *)params;
    uint32_t notificationValue;
    for (;;)
    {
        notificationValue = ulTaskNotifyTake(pdTRUE, portMAX_DELAY);

        if (notificationValue == 1)
        {
            circular_buffer_push(shared_memory->hypersonic_buffer, get_distance(shared_memory->hyper_config));
        }
    }
}
