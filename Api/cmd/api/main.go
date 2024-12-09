package main

import (
	"collector/pkg/config"
	"collector/internal/webApi"
)

func main() {
    apiCfg := &config.ApiConfig{}
	configuration := config.NewConfigurationManager()
	configuration.LoadConfig("./cmd/api/config.json", apiCfg)
    webApi.InitWebApi(*apiCfg)
}