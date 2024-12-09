package routes

import (
	"video-service/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterVideoRoutes(router *gin.Engine, videoController *controllers.VideoController) {
	router.POST("/upload", videoController.UploadVideo)
	router.GET("/metadata/:id", videoController.GetMetadata)
}
