package main

import (
	"collector/internal/webApi"
	"collector/pkg/config"
	
)

func main() {
    apiCfg := &config.ApiConfig{}
	configuration := config.NewConfigurationManager()
	configuration.LoadConfig("./cmd/api/config.json", apiCfg)
    webApi.InitWebApi(*apiCfg)
}