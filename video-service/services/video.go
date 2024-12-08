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

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// VideoService struct to hold DB reference
type VideoService struct {
	DB *mongo.Database
}

// NewVideoService initializes a new VideoService
func NewVideoService(client *mongo.Client) *VideoService {
	return &VideoService{DB: client.Database("video_service")}
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

// UploadToAzureBlob uploads a file to Azure Blob Storage
func (vs *VideoService) UploadToAzureBlob(containerName, fileName string, file io.Reader) (string, error) {
	accountName := utils.GetEnv("AZURE_ACCOUNT_NAME", "")
	accountKey := utils.GetEnv("AZURE_ACCOUNT_KEY", "")

	// Ensure credentials are provided
	if accountName == "" || accountKey == "" {
		return "", errors.New("Azure Blob Storage credentials not provided")
	}

	// Create Azure Blob credentials
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", fmt.Errorf("failed to create Azure Blob credentials: %w", err)
	}

	// Construct the service client
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net", accountName)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create Azure Blob client: %w", err)
	}

	// Create the container client
	containerClient := client.NewContainerClient(containerName)

	// Create the blob client
	blobClient := containerClient.NewBlockBlobClient(fileName)

	// Upload the blob
	_, err = blobClient.Upload(context.Background(), file, nil)
	if err != nil {
		return "", fmt.Errorf("failed to upload blob: %w", err)
	}

	// Return the blob URL
	return blobClient.URL(), nil
}