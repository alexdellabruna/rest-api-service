package main

import (
	"database/sql"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"task-red-hat/internal/handlers"
	"task-red-hat/internal/routes"
	"task-red-hat/internal/storage"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// "/objects/{bucketID}/{objectID}"

var isReady atomic.Bool

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// using SQLite for local storage
	// ignore if the directory already exists
	os.MkdirAll("./db_data", os.ModePerm)

	log.Debug().Msg("Connecting to database...")
	db, err := sql.Open("sqlite3", "./db_data/local.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// create tables if they don't exist
	log.Debug().Msg("Creating tables if they don't exist...")
	db.Exec("CREATE TABLE IF NOT EXISTS buckets (id TEXT PRIMARY KEY)")
	db.Exec("CREATE TABLE IF NOT EXISTS objects (id TEXT PRIMARY KEY, bucket_id TEXT NOT NULL, content TEXT NOT NULL, sha256_hash TEXT NOT NULL, FOREIGN KEY (bucket_id) REFERENCES buckets(id))")

	// get user-specified ip listening address from environment variable
	listeningAddress := os.Getenv("LISTENING_ADDRESS")
	if listeningAddress == "" {
		listeningAddress = "0.0.0.0"
	}

	httpPort := os.Getenv("HTTP_PORT")

	if httpPort == "" {
		httpPort = "8080" // default to 8080 if not set
	}

	// check if the port is a valid integer and within the valid range
	httpPortInt, err := strconv.Atoi(httpPort)
	if err != nil || httpPortInt <= 0 || httpPortInt > 65535 {
		panic("HTTP_PORT must be a valid integer")
	}

	log.Info().Msgf("Starting server on %s:%s", listeningAddress, httpPort)

	genericHTTPHandler := &handlers.GenericHTTPHandler{DBHandler: &storage.DBHandler{DB: db}}

	router := gin.Default()

	routes.RegisterRoutes(router, genericHTTPHandler, &isReady)

	router.Run(listeningAddress + ":" + httpPort)
}
