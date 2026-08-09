package storage

import "database/sql"

type DBHandler struct {
	DB *sql.DB
}

type Bucket struct {
	ID string `json:"id"`
}

type Object struct {
	ID         string `json:"id"`
	BucketID   string `json:"bucket_id"`
	Content    string `json:"content"`
	SHA256Hash string `json:"sha256_hash"`
}
