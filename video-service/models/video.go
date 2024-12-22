package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoMetadata struct {
	ID            primitive.ObjectID    `bson:"_id,omitempty"`      // MongoDB ObjectID
	Title         string    `bson:"title"`             // Video title
	Tags          []string  `bson:"tags"`              // Tags associated with the video
	Duration      int       `bson:"duration"`          // Video duration in seconds
	URL           string    `bson:"url"`               // Video URL
	UploadedAt    time.Time `bson:"uploaded_at"`       // Timestamp of upload
	Thumbnail     string    `bson:"thumbnail"`         // Thumbnail URL
	ThumbnailType string    `bson:"thumbnail_type"`    // Thumbnail type (image or video)
	ContentType   string    `bson:"content_type"`      // Video content type (e.g., video/mp4)
}
