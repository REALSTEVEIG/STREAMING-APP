package routes

import (
	"video-service/controllers"
	"video-service/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterVideoRoutes(router *gin.Engine, db *mongo.Client) {
	videoService := services.NewVideoService(db)
	videoController := controllers.NewVideoController(videoService)

	router.POST("/upload", videoController.UploadVideo)
	router.GET("/metadata/:id", videoController.GetMetadata)
}
