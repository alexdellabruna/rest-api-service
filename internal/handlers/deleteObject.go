package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// DELETE /objects/{bucket}/{objectID}

func (gh *GenericHTTPHandler) DeleteObject(ctx *gin.Context) {
	var req ApiRequestGetPutDelete
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind URI parameters")
		ctx.JSON(http.StatusBadRequest, ResponseTemplatePutDelete{
			Message: "Invalid request",
		})
		return
	}

	bucket := req.Bucket
	objectID := req.ObjectID

	rowsAffected, err := gh.DBHandler.DeleteObjectByID(ctx, bucket, objectID)

	if rowsAffected == 0 {
		log.Warn().Msgf("Object with ID %s in bucket %s not found", objectID, bucket)
		ctx.JSON(http.StatusNotFound, ResponseTemplatePutDelete{
			Bucket:   bucket,
			ObjectID: objectID,
			Message:  "Object not found",
		})
		return
	}

	if err != nil {
		log.Error().Err(err).Msgf("Failed to delete object with ID %s in bucket %s", objectID, bucket)
		ctx.JSON(http.StatusInternalServerError, ResponseTemplatePutDelete{
			Bucket:   bucket,
			ObjectID: objectID,
			Message:  "Failed to delete object",
		})
		return
	}

	log.Info().Msgf("Successfully deleted object with ID %s in bucket %s", objectID, bucket)
	ctx.JSON(http.StatusOK, ResponseTemplatePutDelete{
		Bucket:   bucket,
		ObjectID: objectID,
		Message:  "Object deleted successfully",
	})
}
