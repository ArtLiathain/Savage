package client

import (
	"collector/pkg/collectorsdk"
	"collector/pkg/config"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
)

type DataSnapshot = collectorsdk.DataSnapshot

type ESPPayload struct {
	Count             uint8
	BufferData        []uint8
	SamplingFrequency uint8
	Guid              string
}
type ESPProtocol int

const (
	HealthCheck ESPProtocol = 1
	HyperSonic  ESPProtocol = 2
	Led         ESPProtocol = 3
	Sleep       ESPProtocol = 4
)

var protocolToMetric = map[ESPProtocol]string{
	HyperSonic: "hypersonic_distance",
	Led:        "led_status",
}

// getOSSnapshot collects and logs OS metrics
func getOSSnapshot(channel chan<- DataSnapshot) {
	// Use the logger with predefined fields
	logger := newLogger()

	v, err := mem.VirtualMemory()
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err.Error()}).Error("Error fetching memory stats")
		close(channel)
		return
	}

	metrics := make(map[string]float64)
	metrics["ram_total"] = float64(v.Total / 1024 / 1024)
	metrics["ram_usage"] = float64(v.Used / 1024 / 1024)
	metrics["ram_usage_percent"] = float64(v.UsedPercent)

	info, err := host.Info()
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err.Error()}).Error("Error fetching host info")
		close(channel)
		return
	}

	cpuUsage, err := cpu.Percent(0, false)  // Get CPU usage percentage
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err.Error()}).Error("Error fetching CPU stats")
		close(channel)
		return
	}
	metrics["cpu_usage_percent"] = cpuUsage[0] // We get a slice with one element, so take the first

	// Fetch disk usage stats
	diskStats, err := disk.Usage("/")
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err.Error()}).Error("Error fetching disk stats")
		close(channel)
		return
	}
	metrics["disk_total"] = float64(diskStats.Total / 1024 / 1024 / 1024) // Total disk space in GB
	metrics["disk_used"] = float64(diskStats.Used / 1024 / 1024 / 1024)   // Used disk space in GB
	metrics["disk_usage_percent"] = float64(diskStats.UsedPercent) 

	logger.WithFields(logrus.Fields{"host_id": info.HostID, "hostname": info.Hostname}).Info("Sending OS snapshot")

	channel <- createStandardSnapshot(info.HostID, info.Hostname, metrics, 0)
	close(channel)
}

// createStandardSnapshot creates a standard snapshot for a device
func createStandardSnapshot(device_guid string, device_name string, metrics map[string]float64, recording_offset int) DataSnapshot {
	nowTime := time.Now()
	_, offset := nowTime.Zone()

	return DataSnapshot{
		Timestamp:       nowTime.UTC().Add(-1 * time.Duration(recording_offset) * time.Second),
		TimezoneMinutes: (offset % 3600) / 60,
		DeviceId:        device_guid,
		DeviceName:      device_name,
		Metrics:         metrics,
	}
}

// getEspSnapshot collects ESP data and logs the process
func getEspSnapshot(clientCfg config.ClientConfig, protocol ESPProtocol, channel chan<- DataSnapshot) {
	logger := newLogger()

	parsed_data, err := tcpRequest(clientCfg.ESPDeviceHost, protocol, 5)
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err.Error()}).Error("Error in TCP request")
		close(channel)
		return
	}

	if len(parsed_data) == 0 {
		logger.Warn("Received empty data from ESP")
		close(channel)
		return
	}

	metrics := make(map[string]float64)

	// Decode ESP payload
	decodedData, err := decodeESPDataPayload(parsed_data)
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err.Error()}).Error("Error decoding ESP payload")
		close(channel)
		return
	}

	for index, value := range decodedData.BufferData {
		metrics[protocolToMetric[protocol]] = float64(value)

		channel <- createStandardSnapshot(decodedData.Guid,
			clientCfg.ESPDeviceName, metrics,
			(len(decodedData.BufferData)-index)*int(decodedData.SamplingFrequency))
	}

	logger.WithFields(logrus.Fields{"protocol": fmt.Sprint(protocol)}).Info("ESP snapshot completed")
	close(channel)
}

// tcpRequest sends a TCP request to the ESP device and receives a response
func tcpRequest(serverAddr string, protocol ESPProtocol, amount int) ([]byte, error) {
	logger := newLogger()

	timeoutDuration := 2 * time.Second
	conn, err := net.DialTimeout("tcp", serverAddr, timeoutDuration)
	if err != nil {
		logger.WithFields(logrus.Fields{"server": serverAddr, "error": err.Error()}).Error("Error connecting to server")
		return make([]byte, 0), errors.New("connection error")
	}
	defer conn.Close()

	var payload []byte

	switch protocol {
	case HealthCheck:
		payload = []byte{0x00}
	case HyperSonic:
		payload = []byte{0x01, byte(amount)}
	case Led:
		payload = []byte{0x02, byte(amount)}
	case Sleep:
		payload = []byte{0x03}
	default:
		return make([]byte, 0), errors.New("invalid protocol")
	}

	// Send data
	_, err = conn.Write(payload)
	if err != nil {
		logger.WithFields(logrus.Fields{"protocol": fmt.Sprint(protocol), "error": err.Error()}).Error("Error sending data")
		return make([]byte, 0), err
	}

	logger.WithFields(logrus.Fields{"protocol": fmt.Sprint(protocol)}).Info("Payload sent to server")

	// Receive response
	receivedData := make([]byte, 128)
	_, err = conn.Read(receivedData)
	if err != nil {
		logger.WithFields(logrus.Fields{"protocol": fmt.Sprint(protocol), "error": err.Error()}).Error("Error reading response")
		return make([]byte, 0), err
	}

	return receivedData, nil
}

// decodeESPDataPayload decodes the ESP data payload
func decodeESPDataPayload(payload []byte) (ESPPayload, error) {
	var parsed_payload ESPPayload
	index := 0

	parsed_payload.Count = payload[index]
	index++

	parsed_payload.BufferData = payload[index : index+int(parsed_payload.Count)]
	index += int(parsed_payload.Count)

	parsed_payload.SamplingFrequency = payload[index]
	index++

	if len(payload) < index+8 {
		return parsed_payload, fmt.Errorf("payload too short for offset_minutes")
	}
	parsed_payload.Guid = net.HardwareAddr(payload[index : index+6]).String()
	return parsed_payload, nil
}
