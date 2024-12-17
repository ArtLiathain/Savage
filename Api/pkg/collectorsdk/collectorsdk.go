package collectorsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/blang/semver"
)

type DataSnapshot struct {
	Timestamp       time.Time          `json:"timestamp"`
	TimezoneMinutes int                `json:"timezone_minutes"`
	DeviceId        string             `json:"deviceid"`
	DeviceName      string             `json:"devicename"`
	Metrics         map[string]float64 `json:"metrics"`
}

type Metric struct {
	Name  string
	Value float64
}

var snapshotChannel = make(chan DataSnapshot, 30)

var resendOnce sync.Once

// SendDataSnapshots sends data snapshots to the given HTTP URL. It handles errors and resends data if the request fails.
func SendDataSnapshots(snapshots []DataSnapshot, httpPostUrl string, clientVersion string) (string, error) {
	// JSON encode the snapshots data
	jsonData, err := json.Marshal(snapshots)
	if err != nil {
		return "", fmt.Errorf("JSON encoding error: %w", err)
	}

	// Create a new POST request with the encoded data
	request, err := http.NewRequest("POST", httpPostUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request for %s: %w", httpPostUrl, err)
	}

	// Set necessary headers
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Client-Version", clientVersion)

	// Create an HTTP client to execute the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		// If the HTTP request fails, buffer the data in the snapshot channel and retry sending
		go resendSnapshots(httpPostUrl, clientVersion)
		for _, snapshot := range snapshots {
			snapshotChannel <- snapshot
		}
		return "", fmt.Errorf("HTTP request failed for %s: %w", httpPostUrl, err)
	}
	defer response.Body.Close()

	// Return the status code as the result if the request is successful
	return response.Status, nil
}

// resendSnapshots is responsible for periodically attempting to resend the buffered snapshots.
func resendSnapshots(httpPostUrl string, clientVersion string) {
	resendOnce.Do(func() {
		var snapshots []DataSnapshot

		// Try to resend periodically when there's data in the channel
		for {
			select {
			case snapshot, ok := <-snapshotChannel:
				if !ok {
					// If the channel is closed, exit the loop
					return
				}
				// Collect snapshots from the channel
				snapshots = append(snapshots, snapshot)

			default:
				// Attempt to resend if there are snapshots in the slice
				if len(snapshots) > 0 {
					// Try sending the snapshots when we have collected some
					status, err := SendDataSnapshots(snapshots, httpPostUrl, clientVersion)
					if err == nil {
						// If successful, clear the slice and reset the loop
						fmt.Println("Successfully resent snapshots, status:", status)
						snapshots = nil
						return
					} else {
						// Log the error when resend fails
						fmt.Printf("Error while resending snapshots: %v\n", err)
					}
				}
				// Add a short delay before trying again
				time.Sleep(5 * time.Second)
			}
		}
	})
}

// ParsePostRequest is used to parse an HTTP POST request, verify versions, and decode the request body into data snapshots.
func ParsePostRequest(r *http.Request, serverVersion string) ([]DataSnapshot, error) {
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed, only POST is supported")
	}

	// Parse the client version from the request header
	clientVersion, err := semver.Parse(r.Header.Get("Client-Version"))
	if err != nil {
		return nil, fmt.Errorf("invalid client version in header: %w", err)
	}

	// Parse the server version
	apiVersion, err := semver.Parse(serverVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid server version: %w", err)
	}

	// Check if the client version is compatible with the server version
	if apiVersion.GT(clientVersion) {
		return nil, errors.New("client version too low, please update")
	}

	// Decode the JSON body into data snapshots
	var parsedSnapshots []DataSnapshot
	err = json.NewDecoder(r.Body).Decode(&parsedSnapshots)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON request body: %w", err)
	}

	// Return the parsed snapshots
	return parsedSnapshots, nil
}
