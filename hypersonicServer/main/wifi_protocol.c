#include "esp_wifi.h"
#include "../include/wifi_protocol.h"
#include "lwip/sockets.h"
#include "esp_log.h"
#include "esp_efuse.h"      
#include "esp_system.h"

size_t metric_response(uint8_t amount, shared_memory_t *shared_memory, CircularBuffer *cb_to_use, uint8_t *byte_array)
{
    // Allocate byte array large enough to hold 'amount' ints
    size_t byte_count = circular_buffer_to_byte_array(cb_to_use, amount, byte_array);
    memcpy(byte_array + byte_count, &shared_memory->sampling_frequency, sizeof(uint8_t));
    byte_count += sizeof(uint8_t);
    int64_t mac[6];
    esp_wifi_get_mac(WIFI_IF_STA, &mac);
    memcpy(byte_array + byte_count,&mac , sizeof(int64_t));
    byte_count += sizeof(int64_t);

    return byte_count;
}

void custom_receiver(int client_sock, shared_memory_t *shared_memory, TaskHandle_t *led_task_handle)
{
    uint8_t buffer[1024];
    size_t bytes_received;

    // Receive data from the client
    bytes_received = recv(client_sock, buffer, sizeof(buffer), 0);

    if (bytes_received == 0)
    {
        ESP_LOGE("Receiver", "Failed to receive data.");
    }

    // Ensure there is at least one byte to check
    if (bytes_received > 0)
    {
        // Check the first byte and perform corresponding actions
        if (buffer[0] == 0x00)
        {
            ESP_LOGW("Receiver", "Alive");
            send(client_sock, "Healthy", sizeof("Healthy"), 0);
        }
        else if (buffer[0] == 0x01)
        {
            ESP_LOGW("Receiver", "0x01");
            uint8_t amount = buffer[1];
            if (amount == 0 || amount > 10)
            {
                amount = 5;
            }
            uint8_t byte_array[(amount + 1 * sizeof(uint8_t)) + sizeof(uint8_t) + sizeof(int64_t)];
            size_t byte_count = metric_response(buffer[1], shared_memory, shared_memory->hypersonic_buffer, &byte_array);
            send(client_sock, byte_array, byte_count, 0);
        }
        else if (buffer[0] == 0x02)
        {
            ESP_LOGW("Receiver", "0x02");
            uint8_t amount = buffer[1];
            if (amount == 0 || amount > 10)
            {
                amount = 5;
            }
            uint8_t byte_array[(amount + 1 * sizeof(uint8_t)) + sizeof(uint8_t) + sizeof(int64_t)];
            size_t byte_count = metric_response(buffer[1], shared_memory, shared_memory->led_buffer, &byte_array);
            send(client_sock, byte_array, byte_count, 0);
        }
        else if (buffer[0] == 0x03)
        {
            ESP_LOGW("Receiver", "Sleep");
            send(client_sock, "Sleep", sizeof("Sleep"), 0);
            xTaskNotify(*led_task_handle, 2, eSetBits);
            close(client_sock);
            esp_wifi_set_ps(WIFI_PS_MIN_MODEM);
            vTaskDelay(10000 / portTICK_PERIOD_MS); // Sleep for 10 seconds
            xTaskNotify(*led_task_handle, 3, eSetBits);
            esp_wifi_set_ps(WIFI_PS_NONE);
        }
        else
        {
            ESP_LOGW("Receiver", "No Valid bytes recieved.");
        }
    }
    else
    {
        ESP_LOGW("Receiver", "Received data is empty or too small.");
    }
    ESP_LOGI("Receiver", "Closing client socket");
    close(client_sock);
}