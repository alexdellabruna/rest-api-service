package routes

import (
	"net/http"
	"sync/atomic"
	"task-red-hat/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, dbHandler *handlers.DBHandler, isReady *atomic.Bool) {
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	router.GET("/readyz", func(c *gin.Context) {
		if !isReady.Load() {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}

		// if external db is used, here add a check to see if the database is reachable
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	isReady.Store(true)

	router.GET("/objects/:bucket/:objectID", dbHandler.GetObject)
	router.PUT("/objects/:bucket/:objectID", dbHandler.PutObject)
	router.DELETE("/objects/:bucket/:objectID", dbHandler.DeleteObject)
}
