package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// "/objects/{bucketID}/{objectID}"

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
	db.Exec("CREATE TABLE IF NOT EXISTS objects (id TEXT PRIMARY KEY, bucket_id TEXT NOT NULL, content TEXT NOT NULL, md5_hash TEXT NOT NULL, FOREIGN KEY (bucket_id) REFERENCES buckets(id))")

	// get user-specified port from environment variable
	httpPort := os.Getenv("HTTP_PORT")

	if httpPort == "" {
		httpPort = "8080" // default to 8080 if not set
	}

	// check if the port is a valid integer and within the valid range
	httpPortInt, err := strconv.Atoi(httpPort)
	if err != nil || httpPortInt <= 0 || httpPortInt > 65535 {
		panic("HTTP_PORT must be a valid integer")
	}

	dbHandler := &DBHandler{DB: db}

	router := gin.Default()
	router.GET("/objects/:bucketID/:objectID", dbHandler.getObject)
	router.PUT("/objects/:bucketID/:objectID", dbHandler.putObject)
	router.DELETE("/objects/:bucketID/:objectID", dbHandler.deleteObject)

	// assuming the server is running on all interfaces
	router.Run(":" + httpPort)
}
