package tests

import (
	"database/sql"
	"os"
	"testing"
)

func dbConnectAndInit(t *testing.T) *sql.DB {
	os.MkdirAll("../db_data", os.ModePerm)
	db, err := sql.Open("sqlite3", "../db_data/local.db")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// create tables if they don't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS buckets (id TEXT PRIMARY KEY)")
	if err != nil {
		t.Fatalf("Failed to create buckets table: %v", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS objects (id TEXT PRIMARY KEY, bucket_id TEXT NOT NULL, content TEXT NOT NULL, sha256_hash TEXT NOT NULL, FOREIGN KEY (bucket_id) REFERENCES buckets(id))")
	if err != nil {
		t.Fatalf("Failed to create objects table: %v", err)
	}

	return db
}
