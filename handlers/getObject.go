package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (dbh *DBHandler) GetObject(ctx *gin.Context) {
	// validate input against schema
	var req ApiRequestGetPutDelete
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ResponseTemplateGet{
			ResponseTemplatePutDelete: ResponseTemplatePutDelete{
				Message: "Invalid request",
				Error:   err.Error(),
			},
		})
		return
	}

	bucket := req.Bucket
	objectID := req.ObjectID

	rowsObjectContent, err := dbh.DB.QueryContext(ctx, "SELECT content FROM objects WHERE bucket_id = ? AND id = ?", bucket, objectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseTemplateGet{
			ResponseTemplatePutDelete: ResponseTemplatePutDelete{
				Message: "Failed to get object",
				Error:   err.Error(),
			},
		})
		return
	}
	defer rowsObjectContent.Close()

	if !rowsObjectContent.Next() {
		ctx.JSON(http.StatusNotFound, ResponseTemplateGet{
			ResponseTemplatePutDelete: ResponseTemplatePutDelete{
				Bucket:   "",
				ObjectID: "",
				Message:  "Object not found",
				Error:    "Not found",
			},
		})
		return
	}

	var content string
	err = rowsObjectContent.Scan(&content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseTemplateGet{
			ResponseTemplatePutDelete: ResponseTemplatePutDelete{
				Message: "Failed to read result",
				Error:   err.Error(),
			},
		})
		return
	}
	ctx.JSON(http.StatusOK, ResponseTemplateGet{
		// first two fields are maintened for reference, the content field is the actual object content
		ResponseTemplatePutDelete: ResponseTemplatePutDelete{
			Bucket:   bucket,
			ObjectID: objectID,
			Message:  "Object retrieved successfully",
			Error:    "",
		},
		Content: content,
	})
}
