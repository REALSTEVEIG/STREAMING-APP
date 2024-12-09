package controllers

import (
	"net/http"
	"time"
	"video-service/models"
	"video-service/services"

	"video-service/utils"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	Service *services.VideoService
}

// NewVideoController initializes a new VideoController
func NewVideoController(service *services.VideoService) *VideoController {
	return &VideoController{Service: service}
}

func (vc *VideoController) UploadVideo(c *gin.Context) {
	// Validate title
	title, err := utils.ValidateRequiredField(c, "title", "Title is required")
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate and retrieve file
	file, header, err := utils.ValidateFile(c, "file")
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	// Extract content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Content-Type is required")
		return
	}

	// Upload to S3 and calculate video duration
	fileName := header.Filename
	url, duration, err := vc.Service.ProcessAndUploadVideo(fileName, contentType, file)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Save metadata
	err = vc.Service.SaveVideoMetadata(models.VideoMetadata{
		Title:       title,
		Tags:        c.PostFormArray("tags"),
		Duration:    duration,
		URL:         url,
		UploadedAt:  time.Now(),
		ContentType: contentType,
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to save metadata")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"message":  "Video uploaded successfully",
		"url":      url,
		"metadata": "metadata saved successfully",
	})
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
