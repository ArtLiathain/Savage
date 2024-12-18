package webApi

import (
	"collector/internal/databaseApi"
	"collector/pkg/collectorsdk"
	"collector/pkg/config"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
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

var low_power bool

func (configManager *HandlerConfig) recieveSnapshot(w http.ResponseWriter, r *http.Request) {
	logger := newLogger()

	logger.Info("Received snapshot data")

	parsedSnapshots, err := collectorsdk.ParsePostRequest(r, configManager.config.ApiVersion)
	if err != nil {
		logger.WithError(err).Error("Error parsing post request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error adding metrics %s", err)))
		return
	}

	db, err := sql.Open(configManager.config.DatabaseType, configManager.config.DatabaseName)
	if err != nil {
		logger.WithError(err).Fatal("Error opening database")
		return
	}
	defer db.Close()

	for _, snapshot := range parsedSnapshots {
		err = insertMetrics(db, snapshot)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"device_id": snapshot.DeviceId,
			}).WithError(err).Error("Error inserting metrics")
		} else {
			logger.WithFields(logrus.Fields{
				"device_id": snapshot.DeviceId,
			}).Info("Metrics inserted successfully")
		}
	}
	if low_power {
		low_power = false
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusCreated)

	}
}

func insertMetrics(db *sql.DB, snapshot DataSnapshot) error {
	ctx := context.Background()
	queries := databaseApi.New(db)
	logger := newLogger()
	device_id, err := queries.GetOneDevice(ctx, snapshot.DeviceId)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"device_id": snapshot.DeviceId,
		}).WithError(err).Warn("Device not found, inserting new device")

		device_id, err = queries.InsertDevice(ctx, databaseApi.InsertDeviceParams{DeviceGuid: snapshot.DeviceId, DeviceName: snapshot.DeviceName})
		if err != nil {
			logger.WithError(err).Error("Failed to insert device")
			return err
		}
	}

	nowTime := time.Now()
	_, offset := nowTime.Zone()

	snapshot_id, err := queries.InsertSnapshotTime(ctx, databaseApi.InsertSnapshotTimeParams{
		ClientUtcTime:         snapshot.Timestamp,
		ClientTimezoneMinutes: int64(snapshot.TimezoneMinutes),
		ServerUtcTime:         nowTime,
		ServerTimezoneMinutes: int64(offset),
	})
	if err != nil {
		logger.WithError(err).Error("Failed to insert snapshot time")
		return err
	}

	for metric, value := range snapshot.Metrics {
		metric_id, err := queries.GetOneMetricLookup(ctx, metric)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"metric": metric,
			}).WithError(err).Warn("Metric not found")
			continue
		}

		insert_data := databaseApi.InsertMetricParams{
			SnapshotID: snapshot_id,
			MetricID:   metric_id,
			DeviceID:   device_id,
			Value:      value,
		}

		err = queries.InsertMetric(ctx, insert_data)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"metric": metric,
			}).WithError(err).Error("Failed to insert metric")
			return err
		}

		logger.WithFields(logrus.Fields{
			"metric": metric,
		}).Debug("Metric inserted successfully")
	}
	return nil
}

func (configManager *HandlerConfig) readSnapshots(w http.ResponseWriter, r *http.Request) {
	logger := newLogger()

	db, err := sql.Open(configManager.config.DatabaseType, configManager.config.DatabaseName)
	if err != nil {
		logger.WithError(err).Fatal("Error opening database")
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
	limit, _ := strconv.ParseInt(queryParams.Get("limit"), 10, 64)
	limitNumber := sql.NullInt64{Int64: limit, Valid: limit != 0}
	offset, _ := strconv.ParseInt(queryParams.Get("page"), 10, 64)
	if offset < 2 {
		offset = 1
	}
	offsetNumber := sql.NullInt64{Int64: (offset - 1) * limitNumber.Int64, Valid: true}

	logger.WithFields(logrus.Fields{
		"device_id":   devicefilter,
		"snapshot_id": snapshotfilter,
		"metric_id":   metricfilter,
		"limit":       limitNumber,
		"offset":      offsetNumber,
	}).Debug("Applied filters for reading snapshots")

	filter := databaseApi.GetFilteredMetricsParams{
		DeviceID: devicefilter, MetricID: metricfilter, SnapshotID: snapshotfilter,
		Limit: limitNumber, Offset: offsetNumber,
	}

	response, err := queries.GetFilteredMetrics(ctx, filter)
	if err != nil {
		logger.WithError(err).Error("Error fetching filtered metrics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		logger.WithError(err).Error("Error encoding JSON")
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	logger.Info("Snapshots successfully retrieved and sent")
}

func (configManager *HandlerConfig) resetESP(w http.ResponseWriter, _ *http.Request) {
	logger := newLogger()
	low_power = true
	logger.Info("Low power requst recieved")

	w.WriteHeader(http.StatusCreated)
}

func (configManager *HandlerConfig) getAllDevices(w http.ResponseWriter, _ *http.Request) {
	logger := newLogger()

	db, err := sql.Open(configManager.config.DatabaseType, configManager.config.DatabaseName)
	if err != nil {
		logger.WithError(err).Fatal("Error opening database")
	}
	defer db.Close()

	queries := databaseApi.New(db)
	ctx := context.Background()

	response, err := queries.GetAllDevices(ctx)
	if err != nil {
		logger.WithError(err).Error("Error fetching devices")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		logger.WithError(err).Error("Error encoding JSON")
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	logger.Info("Devices successfully retrieved and sent")
}
func (configManager *HandlerConfig) getAllMetricTypes(w http.ResponseWriter, _ *http.Request) {
	logger := newLogger()

	db, err := sql.Open(configManager.config.DatabaseType, configManager.config.DatabaseName)
	if err != nil {
		logger.WithError(err).Fatal("Error opening database")
	}
	defer db.Close()

	queries := databaseApi.New(db)
	ctx := context.Background()

	response, err := queries.GetAllMetricLookup(ctx)
	if err != nil {
		logger.WithError(err).Error("Error fetching metric lookups")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		logger.WithError(err).Error("Error encoding JSON")
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	logger.Info("Metric Types successfully retrieved and sent")
}
