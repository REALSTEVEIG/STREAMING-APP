package models

import "time"

type VideoMetadata struct {
	ID          string    `bson:"_id,omitempty"`      // MongoDB ObjectID
	Title       string    `bson:"title"`             // Video title
	Tags        []string  `bson:"tags"`              // Tags associated with the video
	Duration    int       `bson:"duration"`          // Video duration in seconds
	URL         string    `bson:"url"`               // Storage URL
	UploadedAt  time.Time `bson:"uploaded_at"`       // Timestamp of upload
	Thumbnail   string    `bson:"thumbnail"`         // Thumbnail URL
	ContentType string    `bson:"content_type"`      // e.g., "video/mp4"
}
