package routes

import (
	"video-service/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterVideoRoutes(router gin.IRouter, videoController *controllers.VideoController) {
	router.POST("/upload", videoController.UploadVideo)
	router.GET("/:id", videoController.GetMetadata)
}
