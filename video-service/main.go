package main

import (
    "context"
    "log"
    "video-service/controllers"
    "video-service/routes"
    "video-service/services"
    "video-service/utils"

    "github.com/gin-gonic/gin"
    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    _"video-service/docs"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// @title Video Service API
// @version 1.0
// @description API for video uploads and metadata management
// @host localhost:8080
// @BasePath /api/videos

func main() {
    utils.LoadEnv()

    clientOptions := options.Client().ApplyURI(utils.GetEnv("MONGO_URI", "mongodb://localhost:27017"))
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    videoService, err := services.NewVideoService(client)
    if err != nil {
        log.Fatalf("Failed to initialize VideoService: %v", err)
    }

    videoController := controllers.NewVideoController(videoService)

    router := gin.Default()

    // Prefix all video routes with /api/video
    apiGroup := router.Group("/api/videos")
    routes.RegisterVideoRoutes(apiGroup, videoController)

    // Swagger docs route
    router.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    log.Println("Starting server on port 8080...")
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
