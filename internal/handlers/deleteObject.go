package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (gh *GenericHTTPHandler) DeleteObject(ctx *gin.Context) {
	var req ApiRequestGetPutDelete
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ResponseTemplatePutDelete{
			Message: "Invalid request",
		})
		return
	}

	bucket := req.Bucket
	objectID := req.ObjectID

	rowsAffected, err := gh.DBHandler.DeleteObjectByID(ctx, bucket, objectID)

	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, ResponseTemplatePutDelete{
			Bucket:   bucket,
			ObjectID: objectID,
			Message:  "Object not found",
		})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseTemplatePutDelete{
			Bucket:   bucket,
			ObjectID: objectID,
			Message:  "Failed to delete object",
		})
		return
	}

	ctx.JSON(http.StatusOK, ResponseTemplatePutDelete{
		Bucket:   bucket,
		ObjectID: objectID,
		Message:  "Object deleted successfully",
	})
}
