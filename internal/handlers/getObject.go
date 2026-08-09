package handlers

import (
	"database/sql"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

// GET /objects/{bucket}/{objectID}

func (gh *GenericHTTPHandler) GetObject(ctx *gin.Context) {
	// validate input against schema
	var req ApiRequestGetPutDelete
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind URI parameters")
		ctx.JSON(http.StatusBadRequest, ResponseTemplateGet{
			ResponseTemplatePutDelete: ResponseTemplatePutDelete{
				Bucket:   req.Bucket,
				ObjectID: req.ObjectID,
				Message:  "Invalid request",
			},
			Content: "",
		})
		return
	}

	bucket := req.Bucket
	objectID := req.ObjectID

	objectObj, err := gh.DBHandler.GetObjectByID(ctx, bucket, objectID)
	if err != nil {

		if err == sql.ErrNoRows {
			log.Warn().Msgf("Object with ID %s in bucket %s not found", objectID, bucket)
			ctx.JSON(http.StatusNotFound, ResponseTemplateGet{
				ResponseTemplatePutDelete: ResponseTemplatePutDelete{
					Bucket:   bucket,
					ObjectID: objectID,
					Message:  "Object not found",
				},
				Content: "",
			})
			return
		}

		log.Error().Err(err).Msgf("Failed to get object with ID %s in bucket %s", objectID, bucket)
		ctx.JSON(http.StatusInternalServerError, ResponseTemplateGet{
			ResponseTemplatePutDelete: ResponseTemplatePutDelete{
				Bucket:   bucket,
				ObjectID: objectID,
				Message:  "Failed to get object",
			},
			Content: "",
		})
		return
	}

	log.Info().Msgf("Successfully retrieved object with ID %s in bucket %s", objectID, bucket)
	ctx.JSON(http.StatusOK, ResponseTemplateGet{
		// first two fields are maintened for reference, the content field is the actual object content
		ResponseTemplatePutDelete: ResponseTemplatePutDelete{
			Bucket:   bucket,
			ObjectID: objectID,
			Message:  "Object retrieved successfully",
		},
		Content: objectObj.Content,
	})
}
