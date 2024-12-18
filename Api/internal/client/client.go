package client

import (
	"collector/pkg/collectorsdk"
	"collector/pkg/config"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var programName string
var debugLevel int

// newLogger creates a new logger with default fields and outputs to both console and a log file
func newLogger() *logrus.Entry {
	// Create a new logger instance
	logger := logrus.New()

	// Set log format to JSON
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Open or create a log file
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		logger.Fatal(err)
	}

	// Use io.MultiWriter to send log entries to both the console and the log file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger.SetOutput(multiWriter)
	logger.SetLevel(logrus.Level(debugLevel))
	// Include default program name field
	logger2 := logger.WithField("program_name", programName)

	return logger2
}

// sendStatistics sends collected statistics and logs the process
func sendStatistics(clientCfg config.ClientConfig) {
	// Initialize the logger with default fields
	logger := newLogger()

	// Log the start of the function with a program tag
	logger.WithFields(logrus.Fields{
		"host":           clientCfg.HostURL,
		"client_version": clientCfg.ClientVersion,
	}).Info("Starting to send statistics")

	snapshots := make([]DataSnapshot, 12)
	osChannel := make(chan DataSnapshot, 1)
	hypersonicChannel := make(chan DataSnapshot, 10)
	ledChannel := make(chan DataSnapshot, 10)

	// Start the go routines
	logger.Debug("Starting go routines for snapshots")
	go getOSSnapshot(osChannel)
	go getEspSnapshot(clientCfg, HyperSonic, hypersonicChannel)
	go getEspSnapshot(clientCfg, Led, ledChannel)

	for snapshot := range osChannel {
		snapshots = append(snapshots, snapshot)
		logger.Debug("Received snapshot from OS channel")
	}

	for snapshot := range hypersonicChannel {
		snapshots = append(snapshots, snapshot)
		logger.Debug("Received snapshot from Hypersonic channel")
	}

	for snapshot := range ledChannel {
		snapshots = append(snapshots, snapshot)
		logger.Debug("Received snapshot from LED channel")
	}

	// Sending the snapshots
	httpposturl := clientCfg.HostURL + "/add"
	status, err := collectorsdk.SendDataSnapshots(snapshots, httpposturl, clientCfg.ClientVersion)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to send snapshots")
		return
	}

	if status[0:3] == "200" {
		logger.WithFields(logrus.Fields{
			"program": programName,
			"status":  status[0:3],
		}).Info("Snapshots successfully sent")
		// Send sleep message to ESP
		tcpRequest(clientCfg.ESPDeviceHost, Sleep, 0)
		logger.Debug("Sent sleep signal to ESP")
	} else if status[0:3] == "201" {
		logger.WithFields(logrus.Fields{
			"program": programName,
			"status":  status[0:3],
		}).Info("Snapshots successfully sent")
	} else {
		logger.WithFields(logrus.Fields{
			"program": programName,
			"status":  status[0:3],
		}).Error("Failed to send snapshots")
	}
}

// SendData starts a routine to send statistics periodically
func SendData(clientCfg config.ClientConfig) {
	programName = clientCfg.ProgramName
	debugLevel = clientCfg.DebugLevel
	logger := newLogger()

	for {
		logger.Debug("Starting sendStatistics routine")
		go sendStatistics(clientCfg)
		time.Sleep(time.Duration(clientCfg.FrequencyInNanoseconds))
	}
}