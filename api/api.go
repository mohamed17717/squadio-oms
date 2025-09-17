package api

import (
	"net/http"
	"oms-services/config"

	"github.com/gin-gonic/gin"
)

// HealthCheck returns a simple health check response
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Service is healthy",
	})
}

// RegisterHealthRoutes registers the health check route
func RegisterHealthRoutes() {
	config.Server.GET("/health", HealthCheck)
}
