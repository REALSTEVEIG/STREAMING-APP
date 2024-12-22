package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"video-service/models"
	"video-service/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type VideoService struct {
	DB        *mongo.Database
	S3Client  *s3.Client
	Bucket    string
	Uploader  *manager.Uploader
}

// NewVideoService initializes a new VideoService
func NewVideoService(client *mongo.Client) (*VideoService, error) {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(utils.GetEnv("AWS_REGION", "")))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(s3Client)

	return &VideoService{
		DB:       client.Database("video_service_meta"),
		S3Client: s3Client,
		Bucket:   utils.GetEnv("AWS_S3_BUCKET", ""),
		Uploader: uploader,
	}, nil
}

func (vs *VideoService) SaveVideoMetadata(metadata models.VideoMetadata) (models.VideoMetadata, error) {
    collection := vs.DB.Collection("videos")
    result, err := collection.InsertOne(context.TODO(), metadata)
    if err != nil {
        return metadata, err
    }

    // Assert that the InsertedID is of type primitive.ObjectID
    if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
        metadata.ID = oid // Assign the ObjectID to the metadata struct
    } else {
        return metadata, fmt.Errorf("failed to cast InsertedID to ObjectID")
    }

    return metadata, nil
}

func (vs *VideoService) CreateAndSaveMetadata(title string, tags []string, duration int, videoURL, thumbnailURL, thumbnailType, contentType string) (models.VideoMetadata, error) {
    metadata := models.VideoMetadata{
        Title:         title,
        Tags:          tags,
        Duration:      duration,
        URL:           videoURL,
        Thumbnail:     thumbnailURL,
        ThumbnailType: thumbnailType,
        UploadedAt:    time.Now(),
        ContentType:   contentType,
    }
	return vs.SaveVideoMetadata(metadata); 
}

// GetVideoMetadata retrieves video metadata by ID
func (vs *VideoService) GetVideoMetadata(id string) (*models.VideoMetadata, error) {
	collection := vs.DB.Collection("videos")

	// Convert the id string to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid video ID format")
	}

	var metadata models.VideoMetadata
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&metadata)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("metadata not found")
		}
		return nil, err
	}

	return &metadata, nil
}

// ProcessAndUploadVideo handles S3 upload and duration calculation
func (vs *VideoService) ProcessAndUploadVideo(fileName, contentType string, file io.Reader) (string, int, error) {
	// Ensure the content type is a video
	if !utils.IsVideoContentType(contentType) {
		return "", 0, fmt.Errorf("unsupported file type: %s", contentType)
	}

	// Upload the video to S3
	result, err := vs.Uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      &vs.Bucket,
		Key:         &fileName,
		Body:        file,
		ContentType: &contentType,
		ACL:         "public-read",
	})
	if err != nil {
		return "", 0, fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Save the file temporarily to calculate duration
	tmpFilePath, err := utils.SaveTemporaryFile(fileName, file)
	if err != nil {
		return result.Location, 0, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFilePath)

	// Calculate video duration using ffprobe
	duration, err := utils.CalculateVideoDuration(tmpFilePath)
	if err != nil {
		return result.Location, 0, fmt.Errorf("failed to calculate video duration: %w", err)
	}

	return result.Location, duration, nil
}

func (vs *VideoService) UploadThumbnail(fileName, contentType string, file io.Reader) (string, error) {
	if !utils.IsVideoContentType(contentType) && !utils.IsImageContentType(contentType) {
		return "", fmt.Errorf("invalid thumbnail type: %s", contentType)
	}

	result, err := vs.Uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      &vs.Bucket,
		Key:         &fileName,
		Body:        file,
		ContentType: &contentType,
		ACL:         "public-read",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload thumbnail: %w", err)
	}

	return result.Location, nil
}
