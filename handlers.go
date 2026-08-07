package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (dbh *DBHandler) getObject(ctx *gin.Context) {
	bucketID := ctx.Param("bucketID")
	objectID := ctx.Param("objectID")

	stmt, err := dbh.DB.PrepareContext(ctx, "SELECT content FROM objects WHERE bucket_id = ? AND id = ?")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	res, err := stmt.QueryContext(ctx, bucketID, objectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		return
	}
	defer res.Close()

	if res.Next() {
		var content string
		err = res.Scan(&content)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read result"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			// first two fields are maintened for reference, the content field is the actual object content
			"bucketID": bucketID,
			"objectID": objectID,
			"content":  content,
		})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
	}
}

func (dbh *DBHandler) putObject(ctx *gin.Context) {
	bucketID := ctx.Param("bucketID")
	objectID := ctx.Param("objectID")
	requestBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}

	stmt, err := dbh.DB.PrepareContext(ctx, "INSERT OR REPLACE INTO objects (id, bucket_id, name, content) VALUES (?, ?, ?, ?)")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, objectID, bucketID, objectID, requestBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"bucketID": bucketID,
		"objectID": objectID,
	})
}

func (dbh *DBHandler) deleteObject(ctx *gin.Context) {
	bucketID := ctx.Param("bucketID")
	objectID := ctx.Param("objectID")
	// Mock response
	ctx.JSON(http.StatusOK, gin.H{
		"bucketID": bucketID,
		"objectID": objectID,
	})
}
