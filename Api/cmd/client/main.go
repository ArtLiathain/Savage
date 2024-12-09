package main

import (
	"collector/internal/client"
	"collector/pkg/config"
)


func main() {
    clientCfg := &config.ClientConfig{}
	configuration := config.NewConfigurationManager()
	configuration.LoadConfig("./cmd/client/config.json", clientCfg)
    client.SendData(*clientCfg)
}