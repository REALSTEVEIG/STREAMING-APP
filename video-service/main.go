package main

import (
	"log"
	"video-service/routes"
	"video-service/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	utils.LoadEnv()

	// Initialize database connection
	db, err := utils.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Set up Gin router
	router := gin.Default()

	// Register video routes
	routes.RegisterVideoRoutes(router, db)

	// Start the server
	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
