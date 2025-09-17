package main

import (
	"log"
	"oms-services/api"
	"oms-services/config"
	"oms-services/models"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDatabase()

	// Auto migrate all models
	err := models.AutoMigrate(db)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Create indexes
	err = models.CreateIndexes(db)
	if err != nil {
		log.Fatal("Index creation failed:", err)
	}

	// Register API routes
	api.RegisterHealthRoutes()
	api.RegisterCatalogRoutes()
	api.RegisterOrderRoutes()

	// Start the Gin server
	gin.SetMode(gin.DebugMode)
	config.Server.Run()
}
