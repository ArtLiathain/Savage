#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "esp_log.h"
#include "esp_system.h"
#include "esp_err.h"
#include "esp_wifi.h"
#include "esp_event.h"
#include "esp_netif.h"
#include "lwip/sockets.h"
#include "lwip/dns.h"
#include "nvs_flash.h"
#include "../include/wifi_protocol.h"
#include "sdkconfig.h"
#include "../include/tcp.h"

#define PORT CONFIG_PORT
#define WIFI_SSID CONFIG_WIFI_SSID      
#define WIFI_PASS CONFIG_WIFI_PASS 

static const char *TAG = "TCP_SERVER";

// Wi-Fi event handler
static void wifi_event_handler(void *arg, esp_event_base_t event_base, int32_t event_id, void *event_data)
{
    if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_START)
    {
        esp_wifi_connect();
    }
    else if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_DISCONNECTED)
    {
        esp_wifi_connect();
        ESP_LOGI(TAG, "Reconnecting to Wi-Fi...");
    }
    else if (event_base == IP_EVENT && event_id == IP_EVENT_STA_GOT_IP)
    {
        ip_event_got_ip_t *event = (ip_event_got_ip_t *)event_data;
        ESP_LOGI(TAG, "Got IP Address: " IPSTR, IP2STR(&event->ip_info.ip));
    }
}

// Initialize Wi-Fi
void wifi_init_sta()
{
    nvs_flash_init();
    esp_netif_init();
    esp_event_loop_create_default();

    esp_netif_t *netif = esp_netif_create_default_wifi_sta();

    wifi_init_config_t cfg = WIFI_INIT_CONFIG_DEFAULT();
    esp_wifi_init(&cfg);

    esp_event_handler_register(WIFI_EVENT, ESP_EVENT_ANY_ID, &wifi_event_handler, NULL);
    esp_event_handler_register(IP_EVENT, IP_EVENT_STA_GOT_IP, &wifi_event_handler, NULL);

    wifi_config_t wifi_config = {
        .sta = {
            .ssid = WIFI_SSID,
            .password = WIFI_PASS,
            .bssid_set = false},
    };

    esp_wifi_set_mode(WIFI_MODE_STA);
    esp_wifi_set_config(ESP_IF_WIFI_STA, &wifi_config);
    esp_wifi_start();
}

// TCP server task
void tcp_server_task(void *pvParameters)
{
    tcp_shared_data_t *tcp_shared_memory = (tcp_shared_data_t *)pvParameters;

    struct sockaddr_in server_addr, client_addr;
    socklen_t addr_len = sizeof(client_addr);
    int sock, client_sock;
    char recv_buf[128];
    int recv_len;

    // Create socket
    sock = socket(AF_INET, SOCK_STREAM, IPPROTO_IP);
    if (sock < 0)
    {
        ESP_LOGE(TAG, "Unable to create socket: errno %d", errno);
        vTaskDelete(NULL);
    }

    // Prepare server address
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(PORT);

    // Bind the socket
    if (bind(sock, (struct sockaddr *)&server_addr, sizeof(server_addr)) < 0)
    {
        ESP_LOGE(TAG, "Unable to bind: errno %d", errno);
        vTaskDelete(NULL);
    }

    // Listen for connections
    if (listen(sock, 1) < 0)
    {
        ESP_LOGE(TAG, "Unable to listen: errno %d", errno);
        vTaskDelete(NULL);
    }

    ESP_LOGI(TAG, "Server listening on port %d", PORT);
    while (1)
    {
        client_sock = accept(sock, (struct sockaddr *)&client_addr, &addr_len);
        if (client_sock < 0)
        {
            ESP_LOGE(TAG, "Unable to accept connection: errno %d", errno);
            break;
        }

        ESP_LOGI(TAG, "Connection accepted");
        custom_receiver(client_sock, tcp_shared_memory->shared_memory, tcp_shared_memory->led_task_handle);
        
        
    }
    // Close the connection
    ESP_LOGI(TAG, "Closing socket");
    close(client_sock);
    close(sock);
}
