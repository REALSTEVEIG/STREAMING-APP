package main

import (
	"context"
	"log"
	"video-service/controllers"
	"video-service/routes"
	"video-service/services"
	"video-service/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	
	utils.LoadEnv()

	// Initialize MongoDB connection
	clientOptions := options.Client().ApplyURI(utils.GetEnv("MONGO_URI", "mongodb://localhost:27017"))
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize VideoService
	videoService, err := services.NewVideoService(client)
	if err != nil {
		log.Fatalf("Failed to initialize VideoService: %v", err)
	}

	// Initialize VideoController
	videoController := controllers.NewVideoController(videoService)

	// Create Gin router
	router := gin.Default()

	// Register routes
	routes.RegisterVideoRoutes(router, videoController)

	// Start the server
	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
