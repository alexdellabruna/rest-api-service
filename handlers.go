package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type responseTemplatePutDelete struct {
	BucketID int64  `json:"bucketID"`
	ObjectID int64  `json:"objectID"`
	Message  string `json:"message"`
	Error    string `json:"error"`
}

type responseTemplateGet struct {
	responseTemplatePutDelete
	Content string `json:"content"`
}

func (dbh *DBHandler) getObject(ctx *gin.Context) {
	// validate input against schema
	var req ApiRequestGetDelete
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseTemplateGet{
			responseTemplatePutDelete: responseTemplatePutDelete{
				Message: "Invalid request",
				Error:   err.Error(),
			},
		})
		return
	}

	bucketID := req.BucketID
	objectID := req.ObjectID

	stmt, err := dbh.DB.PrepareContext(ctx, "SELECT content FROM objects WHERE bucket_id = ? AND id = ?")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseTemplateGet{
			responseTemplatePutDelete: responseTemplatePutDelete{
				Message: "Failed to get object",
				Error:   err.Error(),
			},
		})
		return
	}
	defer stmt.Close()

	res, err := stmt.QueryContext(ctx, bucketID, objectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseTemplateGet{
			responseTemplatePutDelete: responseTemplatePutDelete{
				Message: "Failed to get object",
				Error:   err.Error(),
			},
		})
		return
	}
	defer res.Close()

	if res.Next() {
		var content string
		err = res.Scan(&content)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, responseTemplateGet{
				responseTemplatePutDelete: responseTemplatePutDelete{
					Message: "Failed to read result",
					Error:   err.Error(),
				},
			})
			return
		}
		ctx.JSON(http.StatusOK, responseTemplateGet{
			// first two fields are maintened for reference, the content field is the actual object content
			responseTemplatePutDelete: responseTemplatePutDelete{
				BucketID: bucketID,
				ObjectID: objectID,
				Message:  "Object retrieved successfully",
				Error:    "",
			},
			Content: content,
		})
	} else {
		ctx.JSON(http.StatusNotFound, responseTemplateGet{
			responseTemplatePutDelete: responseTemplatePutDelete{
				BucketID: bucketID,
				ObjectID: objectID,
				Message:  "Object not found",
				Error:    "Object not found",
			},
		})
	}
}

func (dbh *DBHandler) putObject(ctx *gin.Context) {
	var req ApiRequestPut

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseTemplatePutDelete{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	bucketID := req.BucketID
	objectID := req.ObjectID
	content := req.Content

	stmt, err := dbh.DB.PrepareContext(ctx, "INSERT OR REPLACE INTO objects (id, bucket_id, name, content) VALUES (?, ?, ?, ?)")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseTemplatePutDelete{
			BucketID: bucketID,
			ObjectID: objectID,
			Message:  "Failed to store object",
			Error:    err.Error(),
		})
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, objectID, bucketID, objectID, content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseTemplatePutDelete{
			BucketID: bucketID,
			ObjectID: objectID,
			Message:  "Failed to store object",
			Error:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, responseTemplatePutDelete{
		BucketID: bucketID,
		ObjectID: objectID,
		Message:  "Object stored successfully",
		Error:    "",
	})
}

func (dbh *DBHandler) deleteObject(ctx *gin.Context) {
	var req ApiRequestGetDelete
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responseTemplatePutDelete{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	bucketID := req.BucketID
	objectID := req.ObjectID

	stmt, err := dbh.DB.PrepareContext(ctx, "DELETE FROM objects WHERE bucket_id = ? AND id = ?")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseTemplatePutDelete{
			BucketID: bucketID,
			ObjectID: objectID,
			Message:  "Failed to delete object",
			Error:    err.Error(),
		})
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, bucketID, objectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responseTemplatePutDelete{
			BucketID: bucketID,
			ObjectID: objectID,
			Message:  "Failed to delete object",
			Error:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, responseTemplatePutDelete{
		BucketID: bucketID,
		ObjectID: objectID,
		Message:  "Object deleted successfully",
		Error:    "",
	})
}
