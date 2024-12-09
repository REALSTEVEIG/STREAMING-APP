package controllers

import (
	"net/http"
	"time"
	"video-service/models"
	"video-service/services"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	Service *services.VideoService
}

// NewVideoController initializes a new VideoController
func NewVideoController(service *services.VideoService) *VideoController {
	return &VideoController{Service: service}
}

// UploadVideo handles video uploads and saves metadata
func (vc *VideoController) UploadVideo(c *gin.Context) {
	// Retrieve file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		return
	}
	defer file.Close()

	// Upload to AWS S3 via service
	fileName := header.Filename
	url, err := vc.Service.UploadToS3(fileName, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload video"})
		return
	}

	// Save metadata in MongoDB via service
	metadata := models.VideoMetadata{
		Title:       c.PostForm("title"),
		Tags:        c.PostFormArray("tags"),
		Duration:    0, // Add logic to calculate duration
		URL:         url,
		UploadedAt:  time.Now(),
		ContentType: header.Header.Get("Content-Type"),
	}

	err = vc.Service.SaveVideoMetadata(metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video uploaded successfully", "url": url, "metadata": metadata})
}

// GetMetadata retrieves video metadata from MongoDB
func (vc *VideoController) GetMetadata(c *gin.Context) {
	id := c.Param("id")

	metadata, err := vc.Service.GetVideoMetadata(id)
	if err != nil {
		if err.Error() == "metadata not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Metadata not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"metadata": metadata})
}
