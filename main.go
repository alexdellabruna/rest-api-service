package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"

	"task-red-hat/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// "/objects/{bucketID}/{objectID}"

var isReady atomic.Bool

func main() {

	// using SQLite for local storage
	// ignore if the directory already exists
	os.MkdirAll("./db_data", os.ModePerm)
	db, err := sql.Open("sqlite3", "./db_data/local.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create tables if they don't exist
	db.Exec("CREATE TABLE IF NOT EXISTS buckets (id TEXT PRIMARY KEY)")
	db.Exec("CREATE TABLE IF NOT EXISTS objects (id TEXT PRIMARY KEY, bucket_id TEXT NOT NULL, content TEXT NOT NULL, sha256_hash TEXT NOT NULL, FOREIGN KEY (bucket_id) REFERENCES buckets(id))")

	// get user-specified ip listening address from environment variable
	ipListeningAddress := os.Getenv("IP_LISTENING_ADDRESS")
	if ipListeningAddress == "" {
		ipListeningAddress = "0.0.0.0"
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

	dbHandler := &handlers.DBHandler{DB: db}

	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	router.GET("/readyz", func(c *gin.Context) {
		if !isReady.Load() {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}

		// if external db is used, here add a check to see if the database is reachable
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	isReady.Store(true)

	router.GET("/objects/:bucket/:objectID", dbHandler.GetObject)
	router.PUT("/objects/:bucket/:objectID", dbHandler.PutObject)
	router.DELETE("/objects/:bucket/:objectID", dbHandler.DeleteObject)

	// assuming the server is running on all interfaces
	router.Run(ipListeningAddress + ":" + httpPort)
}
