-- name: GetRamRows :many
SELECT
    ram.ram_id,
    ram.timestamp,
    ram.usage,
    ram.total,
    ram.usage_percent,
    devices.device_id,
    devices.device_name
FROM
    ram
    JOIN devices ON ram.device_id = devices.device_id;

-- name: InsertRamRow :exec
INSERT INTO
    ram (device_id, timestamp, usage, total, usage_percent)
VALUES
    (?, ?, ?, ?, ?);

-- name: InsertDeviceIfNotExists :one
INSERT
OR IGNORE INTO devices (device_id, device_name)
VALUES
    (?, ?) RETURNING *;

