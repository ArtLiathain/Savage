CREATE TABLE IF NOT EXISTS
	devices (
		device_id INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
		device_guid VARCHAR(40) UNIQUE NOT NULL,
		device_name TEXT NOT NULL
	);

CREATE TABLE IF NOT EXISTS
	metric_lookup (
		metric_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		name VARCHAR(40) UNIQUE NOT NULL
	);

CREATE TABLE IF NOT EXISTS
	snapshot_time (
		snapshot_id INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
		client_utc_time timestamp NOT NULL,
		client_timezone_minutes INTEGER NOT NULL,
		server_utc_time timestamp NOT NULL,
		server_timezone_minutes INTEGER NOT NULL
	);

CREATE TABLE IF NOT EXISTS
	metrics (
		snapshot_id INTEGER NOT NULL,
		metric_id INT NOT NULL,
		device_id INT NOT NULL,
		value DECIMAL NOT NULL,
		FOREIGN KEY (snapshot_id) REFERENCES snapshot_time (snapshot_id),
		FOREIGN KEY (metric_id) REFERENCES metric_lookup (metric_id),
		FOREIGN KEY (device_id) REFERENCES devices (device_id),
		PRIMARY KEY (snapshot_id, metric_id, device_id)
	);