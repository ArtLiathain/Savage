package client

import (
	"collector/pkg/collectorsdk"
	"collector/pkg/config"
	"time"
)


func sendStatistics( clientCfg config.ClientConfig) {
	snapshot, err := getOSSnapshot() 
	if err != nil {
		panic("OS SNAPSHOT FAILED")
	}
	httpposturl := clientCfg.HostURL + "/add"
	collectorsdk.SendDataSnapshot(snapshot, httpposturl, clientCfg.ClientVersion)
}

func SendData(clientCfg config.ClientConfig) {
	for {
		go sendStatistics(clientCfg)
		time.Sleep(time.Duration(clientCfg.FrequencyInNanoseconds))
	}
}
