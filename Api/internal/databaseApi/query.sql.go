// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package databaseApi

import (
	"context"
	"database/sql"
	"time"
)

const getAllDevices = `-- name: GetAllDevices :many
SELECT
    device_id, device_guid, device_name
FROM
    devices
`

func (q *Queries) GetAllDevices(ctx context.Context) ([]Device, error) {
	rows, err := q.db.QueryContext(ctx, getAllDevices)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Device
	for rows.Next() {
		var i Device
		if err := rows.Scan(&i.DeviceID, &i.DeviceGuid, &i.DeviceName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllMetricLookup = `-- name: GetAllMetricLookup :many
SELECT 
    metric_id, name
FROM 
    metric_lookup
`

func (q *Queries) GetAllMetricLookup(ctx context.Context) ([]MetricLookup, error) {
	rows, err := q.db.QueryContext(ctx, getAllMetricLookup)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MetricLookup
	for rows.Next() {
		var i MetricLookup
		if err := rows.Scan(&i.MetricID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllMetrics = `-- name: GetAllMetrics :many
SELECT 
    snapshot_id,
    metric_id,
    value
FROM 
    metrics
`

type GetAllMetricsRow struct {
	SnapshotID int64
	MetricID   int64
	Value      float64
}

func (q *Queries) GetAllMetrics(ctx context.Context) ([]GetAllMetricsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllMetrics)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllMetricsRow
	for rows.Next() {
		var i GetAllMetricsRow
		if err := rows.Scan(&i.SnapshotID, &i.MetricID, &i.Value); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllSnapshotTimes = `-- name: GetAllSnapshotTimes :many
SELECT 
    client_utc_time,
    client_timezone_minutes,
    server_utc_time,
    server_timezone_minutes
FROM 
    snapshot_time
`

type GetAllSnapshotTimesRow struct {
	ClientUtcTime         time.Time
	ClientTimezoneMinutes int64
	ServerUtcTime         time.Time
	ServerTimezoneMinutes int64
}

func (q *Queries) GetAllSnapshotTimes(ctx context.Context) ([]GetAllSnapshotTimesRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllSnapshotTimes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllSnapshotTimesRow
	for rows.Next() {
		var i GetAllSnapshotTimesRow
		if err := rows.Scan(
			&i.ClientUtcTime,
			&i.ClientTimezoneMinutes,
			&i.ServerUtcTime,
			&i.ServerTimezoneMinutes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFilteredMetrics = `-- name: GetFilteredMetrics :many
SELECT 
    metrics.value,
    metric_lookup.name AS metric_name,
    devices.device_guid,
    devices.device_name,
    snapshot_time.client_utc_time,
    snapshot_time.client_timezone_minutes,
    snapshot_time.server_utc_time,
    snapshot_time.server_timezone_minutes
FROM 
    metrics
JOIN 
    metric_lookup ON metrics.metric_id = metric_lookup.metric_id
JOIN 
    snapshot_time ON metrics.snapshot_id = snapshot_time.snapshot_id
JOIN 
    devices ON devices.device_id = metrics.device_id
WHERE
  (metrics.device_id = ?1 OR ?1 IS NULL)
  AND (metrics.metric_id = ?2 OR ?2 IS NULL)
  AND (metrics.snapshot_id = ?3 OR ?3 IS NULL)
ORDER BY 
  snapshot_time.snapshot_id DESC
LIMIT 
  ?5  
OFFSET 
  ?4
`

type GetFilteredMetricsParams struct {
	DeviceID   sql.NullInt64
	MetricID   sql.NullInt64
	SnapshotID sql.NullInt64
	Offset     sql.NullInt64
	Limit      sql.NullInt64
}

type GetFilteredMetricsRow struct {
	Value                 float64
	MetricName            string
	DeviceGuid            string
	DeviceName            string
	ClientUtcTime         time.Time
	ClientTimezoneMinutes int64
	ServerUtcTime         time.Time
	ServerTimezoneMinutes int64
}

func (q *Queries) GetFilteredMetrics(ctx context.Context, arg GetFilteredMetricsParams) ([]GetFilteredMetricsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFilteredMetrics,
		arg.DeviceID,
		arg.MetricID,
		arg.SnapshotID,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFilteredMetricsRow
	for rows.Next() {
		var i GetFilteredMetricsRow
		if err := rows.Scan(
			&i.Value,
			&i.MetricName,
			&i.DeviceGuid,
			&i.DeviceName,
			&i.ClientUtcTime,
			&i.ClientTimezoneMinutes,
			&i.ServerUtcTime,
			&i.ServerTimezoneMinutes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMetricsWithDetails = `-- name: GetMetricsWithDetails :many
SELECT 
    metrics.value,
    metric_lookup.name AS metric_name,
    devices.device_guid,
    devices.device_name,
    snapshot_time.client_utc_time,
    snapshot_time.client_timezone_minutes,
    snapshot_time.server_utc_time,
    snapshot_time.server_timezone_minutes
FROM 
    metrics
JOIN 
    metric_lookup ON metrics.metric_id = metric_lookup.metric_id
JOIN 
    snapshot_time ON metrics.snapshot_id = snapshot_time.snapshot_id
JOIN 
    devices ON devices.device_id = metrics.device_id
`

type GetMetricsWithDetailsRow struct {
	Value                 float64
	MetricName            string
	DeviceGuid            string
	DeviceName            string
	ClientUtcTime         time.Time
	ClientTimezoneMinutes int64
	ServerUtcTime         time.Time
	ServerTimezoneMinutes int64
}

func (q *Queries) GetMetricsWithDetails(ctx context.Context) ([]GetMetricsWithDetailsRow, error) {
	rows, err := q.db.QueryContext(ctx, getMetricsWithDetails)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMetricsWithDetailsRow
	for rows.Next() {
		var i GetMetricsWithDetailsRow
		if err := rows.Scan(
			&i.Value,
			&i.MetricName,
			&i.DeviceGuid,
			&i.DeviceName,
			&i.ClientUtcTime,
			&i.ClientTimezoneMinutes,
			&i.ServerUtcTime,
			&i.ServerTimezoneMinutes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOneDevice = `-- name: GetOneDevice :one
SELECT
    device_id
FROM
    devices
WHERE
    devices.device_guid == ?
`

func (q *Queries) GetOneDevice(ctx context.Context, deviceGuid string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getOneDevice, deviceGuid)
	var device_id int64
	err := row.Scan(&device_id)
	return device_id, err
}

const getOneMetricLookup = `-- name: GetOneMetricLookup :one
SELECT 
    metric_id
FROM 
    metric_lookup 
WHERE 
    metric_lookup.name == ?
`

func (q *Queries) GetOneMetricLookup(ctx context.Context, name string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getOneMetricLookup, name)
	var metric_id int64
	err := row.Scan(&metric_id)
	return metric_id, err
}

const insertDevice = `-- name: InsertDevice :one
INSERT OR IGNORE INTO devices (device_guid, device_name)
VALUES (?, ?) RETURNING device_id
`

type InsertDeviceParams struct {
	DeviceGuid string
	DeviceName string
}

func (q *Queries) InsertDevice(ctx context.Context, arg InsertDeviceParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, insertDevice, arg.DeviceGuid, arg.DeviceName)
	var device_id int64
	err := row.Scan(&device_id)
	return device_id, err
}

const insertMetric = `-- name: InsertMetric :exec
INSERT INTO metrics (snapshot_id, metric_id, device_id, value)
VALUES (?, ?, ?, ?)
`

type InsertMetricParams struct {
	SnapshotID int64
	MetricID   int64
	DeviceID   int64
	Value      float64
}

func (q *Queries) InsertMetric(ctx context.Context, arg InsertMetricParams) error {
	_, err := q.db.ExecContext(ctx, insertMetric,
		arg.SnapshotID,
		arg.MetricID,
		arg.DeviceID,
		arg.Value,
	)
	return err
}

const insertMetricLookup = `-- name: InsertMetricLookup :exec
INSERT OR IGNORE INTO metric_lookup (name)
VALUES (?)
`

func (q *Queries) InsertMetricLookup(ctx context.Context, name string) error {
	_, err := q.db.ExecContext(ctx, insertMetricLookup, name)
	return err
}

const insertSnapshotTime = `-- name: InsertSnapshotTime :one
INSERT INTO snapshot_time (client_utc_time, client_timezone_minutes, server_utc_time, server_timezone_minutes)
VALUES (?, ?, ?, ?)
RETURNING snapshot_id
`

type InsertSnapshotTimeParams struct {
	ClientUtcTime         time.Time
	ClientTimezoneMinutes int64
	ServerUtcTime         time.Time
	ServerTimezoneMinutes int64
}

func (q *Queries) InsertSnapshotTime(ctx context.Context, arg InsertSnapshotTimeParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, insertSnapshotTime,
		arg.ClientUtcTime,
		arg.ClientTimezoneMinutes,
		arg.ServerUtcTime,
		arg.ServerTimezoneMinutes,
	)
	var snapshot_id int64
	err := row.Scan(&snapshot_id)
	return snapshot_id, err
}
