package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"video-service/models"
	"video-service/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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
		DB:       client.Database("video_service"),
		S3Client: s3Client,
		Bucket:   utils.GetEnv("AWS_S3_BUCKET", ""),
		Uploader: uploader,
	}, nil
}

// SaveVideoMetadata saves video metadata to MongoDB
func (vs *VideoService) SaveVideoMetadata(metadata models.VideoMetadata) error {
	collection := vs.DB.Collection("videos")
	_, err := collection.InsertOne(context.TODO(), metadata)
	return err
}

// GetVideoMetadata retrieves video metadata by ID
func (vs *VideoService) GetVideoMetadata(id string) (*models.VideoMetadata, error) {
	collection := vs.DB.Collection("videos")

	var metadata models.VideoMetadata
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&metadata)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("metadata not found")
		}
		return nil, err
	}

	return &metadata, nil
}

// UploadToS3 uploads a file to AWS S3
func (vs *VideoService) UploadToS3(fileName string, file io.Reader) (string, error) {
	// Upload file to S3
	result, err := vs.Uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(vs.Bucket),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String("video/mp4"), // Adjust as needed
		ACL:         types.ObjectCannedACLPublicRead, // Optional: Use private for restricted access
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Return the public URL of the uploaded file
	return result.Location, nil
}
