package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (dbh *DBHandler) DeleteObject(ctx *gin.Context) {
	var req ApiRequestGetPutDelete
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ResponseTemplatePutDelete{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	bucket := req.Bucket
	objectID := req.ObjectID

	rowsDeletedObject, err := dbh.DB.ExecContext(ctx, "DELETE FROM objects WHERE bucket_id = ? AND id = ?", bucket, objectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseTemplatePutDelete{
			Bucket:   bucket,
			ObjectID: objectID,
			Message:  "Failed to delete object",
			Error:    err.Error(),
		})
		return
	}

	rowsAffected, err := rowsDeletedObject.RowsAffected()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseTemplatePutDelete{
			Bucket:   bucket,
			ObjectID: objectID,
			Message:  "Failed to determine if object was deleted",
			Error:    err.Error(),
		})
		return
	}

	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, ResponseTemplatePutDelete{
			Bucket:   bucket,
			ObjectID: objectID,
			Message:  "Object not found",
			Error:    "",
		})
		return
	}

	ctx.JSON(http.StatusOK, ResponseTemplatePutDelete{
		Bucket:   bucket,
		ObjectID: objectID,
		Message:  "Object deleted successfully",
		Error:    "",
	})
}
