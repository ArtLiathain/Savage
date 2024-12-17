package webApi

import (
	"collector/internal/databaseApi"
	"collector/pkg/config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

func InitWebApi(apiCfg config.ApiConfig) {

	db, err := sql.Open(apiCfg.DatabaseType, apiCfg.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlBytes, err := os.ReadFile(apiCfg.SchemaPath)
	if err != nil {
		log.Fatal(err)
	}

	sqlStatements := string(sqlBytes)

	_, err = db.Exec(sqlStatements)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	queries := databaseApi.New(db)
	for _, metric := range apiCfg.MetricsLookup {
		queries.InsertMetricLookup(ctx, metric)
	}

	log.Println("Database and tables created successfully")

	handlers := &HandlerConfig{config: apiCfg}
	http.HandleFunc("/read", handlers.readSnapshots)
	http.HandleFunc("/add", handlers.recieveSnapshot)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins, replace with your frontend URL in production
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
	})
	handler := c.Handler(http.DefaultServeMux)
	log.Printf("Listening on %s:%d", apiCfg.Host, apiCfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", apiCfg.Port), handler))
}
