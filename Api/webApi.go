package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"collector/databaseApi"
	"context"
)

type Metric struct {
	Name  string
	Value float64
}

type ReadResponse struct {
	DeviceId   string
	DeviceName string
	Timestamp time.Time
	Entries []Metric
}

type DataSnapshot struct {
    Timestamp  time.Time          `json:"timestamp"`
    MetricType string             `json:"metrictype"`  // Fixed case: from 'Metrictype' to 'MetricType'
    DeviceId   string             `json:"deviceid"`
    DeviceName string             `json:"devicename"`
    Metrics    map[string]float64 `json:"metrics"`
}





func recieveSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var parsedSnapshot DataSnapshot
	err := json.NewDecoder(r.Body).Decode(&parsedSnapshot)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	db, err := sql.Open("sqlite3", "metrics.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert metrics into the table
	err = insertMetrics(db, parsedSnapshot.MetricType, parsedSnapshot)
	if err != nil {
		log.Fatalf("Error inserting metrics: %v", err)
	}

	fmt.Println("Metrics inserted successfully!")
}

// insertMetrics inserts metrics into the specified table
func insertMetrics(db *sql.DB, tableName string, snapshot DataSnapshot) error {
    // Start building the SQL statement dynamically
    columns := []string{"deviceid", "devicename", "timestamp"} // Always insert these columns
    values := []interface{}{snapshot.DeviceId, snapshot.DeviceName, snapshot.Timestamp}

    // Collect the metric column names and values from the snapshot.Metrics map
    for key, value := range snapshot.Metrics {
        columns = append(columns, key)  // Add the metric name (key) as the column name
        values = append(values, value) // Add the metric value as the value
    }

    // Create a placeholder for the SQL statement (?)
    placeholders := make([]string, len(columns))
    for i := range placeholders {
        placeholders[i] = "?" // Placeholder for each column
    }

    // Construct the INSERT SQL query dynamically
    insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, 
        join(columns, ","), join(placeholders, ","))	
	fmt.Println(insertSQL)
    // Prepare the SQL statement for execution
    stmt, err := db.Prepare(insertSQL)
    if err != nil {
        return fmt.Errorf("error preparing query: %v", err)
    }
    defer stmt.Close()
	fmt.Println(values...)
    // Execute the INSERT statement with the dynamic values
    _, err = stmt.Exec(values...)
    if err != nil {
        return fmt.Errorf("error executing query: %v", err)
    }

    return nil
}


func join(items []string, separator string) string {
	result := ""
	for i, item := range items {
		if i > 0 {
			result += separator
		}
		result += item
	}
	return result
}

func readSnapshots(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "metrics.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// metricType := r.URL.Query().Get("metricType")
	queries :=databaseApi.New(db)
	ctx := context.Background()

	response, err:= queries.GetRamRows(ctx)
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

func fetchReadResponses(rows *sql.Rows) ([]ReadResponse, error) {
    var responses []ReadResponse

    // Get column names
    cols, err := rows.Columns()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch column names: %w", err)
    }

    // Process rows
    for rows.Next() {
        // Prepare a slice to hold all values
        values := make([]interface{}, len(cols))
        valuePtrs := make([]interface{}, len(cols))

        for i := range values {
            valuePtrs[i] = &values[i]
        }

        // Scan the row into value pointers
        if err := rows.Scan(valuePtrs...); err != nil {
            return nil, fmt.Errorf("failed to scan row: %w", err)
        }

        // Build the ReadResponse
        response := ReadResponse{
            DeviceId:   values[1].(string),
            DeviceName: values[2].(string),
            Timestamp:  values[3].(time.Time),
        }

        // Collect metrics (columns 4 onward)
        for i := 4; i < len(cols); i++ {
            metric := Metric{
                Name:  cols[i],
                Value: values[i].(float64),
            }
            response.Entries = append(response.Entries, metric)
        }

        responses = append(responses, response)
    }

    return responses, nil
}



func initWebApi() {
	http.HandleFunc("/", readSnapshots)
	http.HandleFunc("/add", recieveSnapshot)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins, replace with your frontend URL in production
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
	})
	handler := c.Handler(http.DefaultServeMux)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
