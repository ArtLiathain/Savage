package collectorsdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type DataSnapshot struct {
	Timestamp       time.Time          `json:"timestamp"`
	TimezoneMinutes int                `json:"timezone_minutes"`
	DeviceId        string             `json:"deviceid"`
	DeviceName      string             `json:"devicename"`
	Metrics         map[string]float64 `json:"metrics"`
}

type Metric struct {
	Name  string
	Value float64
}



func SendDataSnapshot(snapshot DataSnapshot, httpPostUrl string, clientVersion string ) {

	jsonData, err := json.Marshal(snapshot)
	if err != nil {
		log.Fatal("JSON encoding error:", err)
	}

	request, err := http.NewRequest("POST", httpPostUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Request creation error:", err)
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Client-Version", clientVersion)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("HTTP request error:", err)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
}