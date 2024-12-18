set(COMPONENT_REQUIRES bt)

idf_component_register(SRCS "led.c" "main.c" "hypersonic.c" "tcp.c" "circular_buffer.c" "timer.c" "wifi_protocol.c"
                    INCLUDE_DIRS ".")