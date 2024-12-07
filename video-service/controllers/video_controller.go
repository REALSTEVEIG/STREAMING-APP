package controllers

import (
	"net/http"
	"video-service/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type VideoController struct {
	DB *mongo.Client
}

func NewVideoController(db *mongo.Client) *VideoController {
	return &VideoController{DB: db}
}

// UploadVideo handles video upload requests
func (vc *VideoController) UploadVideo(c *gin.Context) {
	// Retrieve file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		return
	}
	defer file.Close()

	// Upload to Azure Blob
	containerName := "videos"
	fileName := header.Filename
	url, err := services.UploadToAzureBlob(containerName, fileName, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload video"})
		return
	}

	// Save metadata in MongoDB (optional)

	c.JSON(http.StatusOK, gin.H{"message": "Video uploaded successfully", "url": url})
}

// GetMetadata retrieves video metadata
func (vc *VideoController) GetMetadata(c *gin.Context) {
	// Placeholder implementation for metadata retrieval
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "metadata": "Sample metadata"})
}
