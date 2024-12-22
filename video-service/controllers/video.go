package controllers

import (
	"mime/multipart"
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

// @Summary Upload a video
// @Description Uploads a video and optional thumbnail to S3 and saves metadata
// @Tags videos
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Video title"
// @Param tags formData []string false "Video tags"
// @Param file formData file true "Video file"
// @Param thumbnail formData file false "Thumbnail (video or image)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /upload [post]
func (vc *VideoController) UploadVideo(c *gin.Context) {
	title, err := utils.ValidateRequiredField(c, "title", "Title is required")
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	videoFile, videoHeader, err := utils.ValidateFile(c, "file")
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	defer videoFile.Close()

	// Validate thumbnail file (optional)
	var thumbnailFile multipart.File
	var thumbnailHeader *multipart.FileHeader
	thumbnailType := ""
	thumbnailURL := ""

	if c.Request.MultipartForm != nil {
		thumbnailFile, thumbnailHeader, _ = c.Request.FormFile("thumbnail")
		if thumbnailFile != nil {
			defer thumbnailFile.Close()
			thumbnailType = thumbnailHeader.Header.Get("Content-Type")

			// Validate thumbnail type
			if !utils.IsVideoContentType(thumbnailType) && !utils.IsImageContentType(thumbnailType) {
				utils.RespondWithError(c, http.StatusBadRequest, "Invalid thumbnail type: must be an image or a video")
				return
			}
		}
	}

	// Process and upload video
	videoURL, duration, err := vc.Service.ProcessAndUploadVideo(videoHeader.Filename, videoHeader.Header.Get("Content-Type"), videoFile)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Process and upload thumbnail if provided
	if thumbnailFile != nil {
		thumbnailURL, err = vc.Service.UploadThumbnail(thumbnailHeader.Filename, thumbnailType, thumbnailFile)
		if err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to upload thumbnail")
			return
		}
	}

	// Save metadata
	err = vc.Service.SaveVideoMetadata(models.VideoMetadata{
		Title:         title,
		Tags:          c.PostFormArray("tags"),
		Duration:      duration,
		URL:           videoURL,
		Thumbnail:     thumbnailURL,
		ThumbnailType: thumbnailType,
		UploadedAt:    time.Now(),
		ContentType:   videoHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to save metadata")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"message":        "Video uploaded successfully",
		"url":            videoURL,
		"thumbnail_url":  thumbnailURL,
		"metadata":       "metadata saved successfully",
	})
}

// @Summary Get video metadata
// @Description Retrieves video metadata by ID
// @Tags videos
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /{id} [get]
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
