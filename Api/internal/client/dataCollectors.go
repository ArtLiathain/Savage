package client

import (
	"collector/pkg/collectorsdk"
	"log"
	"time"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type DataSnapshot = collectorsdk.DataSnapshot

func getOSSnapshot() (DataSnapshot, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return DataSnapshot{}, err
	}
	metrics := make(map[string]float64)
	metrics["ram_total"] = float64(v.Total / 1024 / 1024)
	metrics["ram_usage"] = float64(v.Used / 1024 / 1024)
	metrics["ram_usage_percent"] = float64(v.UsedPercent)
	info, err := host.Info()
    if err != nil {
        log.Fatal(err)
    }


	return createStandardSnapshot(info.HostID, info.Hostname, metrics), nil

}

func createStandardSnapshot(device_guid string, device_name string, metrics map[string]float64) DataSnapshot {
	nowTime := time.Now()
	_, offset := nowTime.Zone()

	return DataSnapshot{Timestamp: nowTime.UTC(), TimezoneMinutes: (offset % 3600) / 60, DeviceId: device_guid, DeviceName: device_name, Metrics: metrics}
}