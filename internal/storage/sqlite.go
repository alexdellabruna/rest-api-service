package storage

import (
	"context"
	"database/sql"
)

func (dbh *DBHandler) GetBucketByID(ctx context.Context, bucketID string) (*Bucket, error) {
	rows, err := dbh.DB.QueryContext(ctx, "SELECT * FROM buckets WHERE id = ?", bucketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var bucket Bucket
		err := rows.Scan(&bucket.ID)
		if err != nil {
			return nil, err
		}
		return &bucket, nil
	}

	return nil, sql.ErrNoRows
}

func (dbh *DBHandler) InsertBucket(ctx context.Context, bucketID string) error {
	_, err := dbh.DB.ExecContext(ctx, "INSERT INTO buckets (id) VALUES (?)", bucketID)
	return err
}

func (dbh *DBHandler) CheckObjectExistance(ctx context.Context, bucketID, objectID, sha256Hash string) (bool, error) {
	rows, err := dbh.DB.QueryContext(ctx, "SELECT id FROM objects WHERE (id = ? OR sha256_hash = ?) AND bucket_id = ?", objectID, sha256Hash, bucketID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (dbh *DBHandler) InsertObject(ctx context.Context, objectID, bucketID, content, sha256Hash string) error {
	_, err := dbh.DB.ExecContext(ctx, "INSERT INTO objects (id, bucket_id, content, sha256_hash) VALUES (?, ?, ?, ?)", objectID, bucketID, content, sha256Hash)
	return err
}

func (dbh *DBHandler) GetObjectByID(ctx context.Context, bucketID, objectID string) (*Object, error) {
	rows, err := dbh.DB.QueryContext(ctx, "SELECT * FROM objects WHERE bucket_id = ? AND id = ?", bucketID, objectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var object Object
		err := rows.Scan(&object.ID, &object.BucketID, &object.Content, &object.SHA256Hash)
		if err != nil {
			return nil, err
		}
		return &object, nil
	}

	return nil, sql.ErrNoRows
}

func (dbh *DBHandler) DeleteObjectByID(ctx context.Context, bucketID, objectID string) (int64, error) {
	result, err := dbh.DB.ExecContext(ctx, "DELETE FROM objects WHERE bucket_id = ? AND id = ?", bucketID, objectID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
