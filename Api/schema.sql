CREATE TABLE
	devices (
		device_id INTEGER PRIMARY KEY NOT NULL,
		device_name TEXT NOT NULL
	);

CREATE TABLE
	ram (
		ram_id INTEGER PRIMARY KEY AUTOINCREMENT,
		device_id INT NOT NULL,
		timestamp DATETIME NOT NULL DEFAULT  CURRENT_TIMESTAMP,
		usage REAL NOT NULL,
		total REAL NOT NULL,
		usage_percent REAL NOT NULL,
		FOREIGN KEY (device_id) REFERENCES devices (device_id)
	);