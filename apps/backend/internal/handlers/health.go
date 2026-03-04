package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health godoc
// GET /health
// Returns server health status and timestamp.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "upsync-api",
	})
}
