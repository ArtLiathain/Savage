-- name: InsertDevice :one
INSERT OR IGNORE INTO devices (device_guid, device_name)
VALUES (?, ?) RETURNING device_id;

-- name: InsertMetricLookup :exec
INSERT OR IGNORE INTO metric_lookup (name)
VALUES (?);

-- name: InsertSnapshotTime :one
INSERT INTO snapshot_time (client_utc_time, client_timezone_minutes, server_utc_time, server_timezone_minutes)
VALUES (?, ?, ?, ?)
RETURNING snapshot_id;

-- name: InsertMetric :exec
INSERT INTO metrics (snapshot_id, metric_id, device_id, value)
VALUES (?, ?, ?, ?);

-- name: GetAllMetricLookup :many
SELECT 
    name
FROM 
    metric_lookup;
-- name: GetOneMetricLookup :one
SELECT 
    metric_id
FROM 
    metric_lookup 
WHERE 
    metric_lookup.name == ?;

-- name: GetAllSnapshotTimes :many
SELECT 
    client_utc_time,
    client_timezone_minutes,
    server_utc_time,
    server_timezone_minutes
FROM 
    snapshot_time;

-- name: GetAllMetrics :many
SELECT 
    snapshot_id,
    metric_id,
    value
FROM 
    metrics;

-- name: GetOneDevice :one
SELECT
    device_id
FROM
    devices
WHERE
    devices.device_guid == ?;

-- name: GetMetricsWithDetails :many
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
    devices ON devices.device_id = metrics.device_id;

-- name: GetFilteredMetrics :many
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
  (metrics.device_id = sqlc.narg('device_id') OR sqlc.narg('device_id') IS NULL)
  AND (metrics.metric_id = sqlc.narg('metric_id') OR sqlc.narg('metric_id') IS NULL)
  AND (metrics.snapshot_id = sqlc.narg('snapshot_id') OR sqlc.narg('snapshot_id') IS NULL);