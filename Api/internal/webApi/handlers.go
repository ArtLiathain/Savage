package webApi

import (
	"collector/internal/databaseApi"
	"collector/pkg/collectorsdk"
	"collector/pkg/config"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type HandlerConfig struct {
	config config.ApiConfig
}

type DataSnapshot = collectorsdk.DataSnapshot
type Metric = collectorsdk.Metric
type ReadResponse struct {
	DeviceId   string
	DeviceName string
	Timestamp  time.Time
	Entries    []Metric
}

func (configManager *HandlerConfig) recieveSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var parsedSnapshot DataSnapshot
	err := json.NewDecoder(r.Body).Decode(&parsedSnapshot)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}
	db, err := sql.Open(configManager.config.DatabaseType, configManager.config.DatabaseName)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	err = insertMetrics(db, parsedSnapshot)
	if err != nil {
		log.Printf("Error inserting metrics: %v", err)
	}

	fmt.Println("Metrics inserted successfully!")
}

func insertMetrics(db *sql.DB, snapshot DataSnapshot) error {
	ctx := context.Background()
	queries := databaseApi.New(db)
	device_id, err := queries.GetOneDevice(ctx, snapshot.DeviceId)
	if err != nil {
		device_id, err = queries.InsertDevice(ctx, databaseApi.InsertDeviceParams{DeviceGuid: snapshot.DeviceId, DeviceName: snapshot.DeviceName})
		if err != nil {
			fmt.Println("HIGE ERROR")
			return err
		}
	}
	nowTime := time.Now()
	_, offset := nowTime.Zone()

	snapshot_id, err := queries.InsertSnapshotTime(ctx,
		databaseApi.InsertSnapshotTimeParams{
			ClientUtcTime:         snapshot.Timestamp,
			ClientTimezoneMinutes: int64(snapshot.TimezoneMinutes),
			ServerUtcTime:         nowTime,
			ServerTimezoneMinutes: int64(offset)})
	if err != nil {
		fmt.Println("HIGE ERROR 2")
		return err
	}

	for metric, value := range snapshot.Metrics {

		metric_id, err := queries.GetOneMetricLookup(ctx, metric)
		if err != nil {
			fmt.Println(err)
			continue
		}
		insert_data := databaseApi.InsertMetricParams{
			SnapshotID: snapshot_id,
			MetricID:   metric_id,
			DeviceID:   device_id,
			Value:      value}
		err = queries.InsertMetric(ctx, insert_data)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil

}

func (configManager *HandlerConfig) readSnapshots(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(configManager.config.DatabaseType, configManager.config.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	queries := databaseApi.New(db)
	ctx := context.Background()
	queryParams := r.URL.Query()

	device_id, _ := strconv.ParseInt(queryParams.Get("device_id"), 10, 64)
	devicefilter := sql.NullInt64{Int64: device_id, Valid: device_id != 0}
	snapshot_id, _ := strconv.ParseInt(queryParams.Get("snapshot_id"), 10, 64)
	snapshotfilter := sql.NullInt64{Int64: snapshot_id, Valid: snapshot_id != 0}
	metric_id, _ := strconv.ParseInt(queryParams.Get("metric_id"), 10, 64)
	metricfilter := sql.NullInt64{Int64: metric_id, Valid: metric_id != 0}

	filter := databaseApi.GetFilteredMetricsParams{DeviceID: devicefilter, MetricID: metricfilter, SnapshotID: snapshotfilter}
	response, err := queries.GetFilteredMetrics(ctx, filter)
	if err != nil {
		log.Fatal(w, err)
		return
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	// Set the response header to JSON and write the JSON data
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
