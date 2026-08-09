package tests

import (
	"database/sql"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func dbConnectAndInit(t *testing.T) *sql.DB {
	os.MkdirAll("./db_data", os.ModePerm)
	db, err := sql.Open("sqlite3", "./db_data/local.db")
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

func seedTestBucketData(t *testing.T, db *sql.DB, ctx *gin.Context, bucket string) {
	t.Helper()
	_, err := db.ExecContext(ctx, "INSERT INTO buckets (id) VALUES (?)", bucket)
	if err != nil {
		t.Errorf("Failed to seed test data: %v", err)
	}
}

func seedTestObjectData(t *testing.T, db *sql.DB, ctx *gin.Context, bucket string, objectID string) {
	t.Helper()
	_, err := db.ExecContext(ctx, "INSERT INTO objects (id, bucket_id, content, sha256_hash) VALUES (?, ?, ?, ?)", objectID, bucket, "This is a test object content.", "abcd1234efgh5678ijkl9012mnop3456qrst7890uvwx1234yzab5678cdef9012")

	if err != nil {
		t.Errorf("Failed to seed test data: %v", err)
	}
}

func cleanupTestBucketData(t *testing.T, db *sql.DB, ctx *gin.Context, bucket string) {
	t.Helper()
	_, err := db.ExecContext(ctx, "DELETE FROM buckets WHERE id = ?", bucket)
	if err != nil {
		t.Errorf("Failed to clean up test data: %v", err)
	}
}

func cleanupTestObjectData(t *testing.T, db *sql.DB, ctx *gin.Context, bucket string, objectID string) {
	t.Helper()
	_, err := db.ExecContext(ctx, "DELETE FROM objects WHERE id = ? AND bucket_id = ?", objectID, bucket)
	if err != nil {
		t.Errorf("Failed to clean up test data: %v", err)
	}
}
