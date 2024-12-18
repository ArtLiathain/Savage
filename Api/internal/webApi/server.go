package webApi

import (
	"collector/internal/databaseApi"
	"collector/pkg/config"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

var programName string
var debugLevel int

func newLogger() *logrus.Entry {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		logger.Fatal(err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger.SetOutput(multiWriter)
	logger.SetLevel(logrus.Level(debugLevel))

	logger2 := logger.WithField("program_name", programName)

	return logger2
}

func InitWebApi(apiCfg config.ApiConfig) {
	programName = apiCfg.ProgramName
	debugLevel = apiCfg.DebugLevel
	logger := newLogger()

	logger.Info("Initializing Web API")

	db, err := sql.Open(apiCfg.DatabaseType, apiCfg.DatabaseName)
	if err != nil {
		logger.WithError(err).Fatal("Failed to open database")
	}
	defer db.Close()
	logger.Info("Database connection established")

	sqlBytes, err := os.ReadFile(apiCfg.SchemaPath)
	if err != nil {
		logger.WithError(err).Fatal("Failed to read schema file")
	}
	sqlStatements := string(sqlBytes)

	_, err = db.Exec(sqlStatements)
	if err != nil {
		logger.WithError(err).Fatal("Failed to execute schema")
	}
	logger.Info("Database schema applied successfully")

	ctx := context.Background()
	queries := databaseApi.New(db)
	for _, metric := range apiCfg.MetricsLookup {
		err := queries.InsertMetricLookup(ctx, metric)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"metric": metric,
			}).WithError(err).Error("Failed to insert metric lookup")
		} else {
			logger.WithFields(logrus.Fields{
				"metric": metric,
			}).Info("Metric lookup inserted successfully")
		}
	}

	handlers := &HandlerConfig{config: apiCfg}
	http.HandleFunc("/read", handlers.readSnapshots)
	http.HandleFunc("/add", handlers.recieveSnapshot)
	http.HandleFunc("/reset", handlers.resetESP)
	http.HandleFunc("/devices", handlers.getAllDevices)
	http.HandleFunc("/metrictypes", handlers.getAllMetricTypes)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins, replace with your frontend URL in production
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
	})
	handler := c.Handler(http.DefaultServeMux)

	logger.Infof("Listening on %s:%d", apiCfg.Host, apiCfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", apiCfg.Port), handler)
	if err != nil {
		logger.WithError(err).Fatal("Server failed to start")
	}
}
