package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"video-service/utils"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// UploadToAzureBlob uploads a file to Azure Blob Storage.
func UploadToAzureBlob(containerName, fileName string, file io.Reader) (string, error) {
	accountName := utils.GetEnv("AZURE_ACCOUNT_NAME", "")
	accountKey := utils.GetEnv("AZURE_ACCOUNT_KEY", "")

	// Create credential
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatalf("Failed to create Azure Blob credentials: %v", err)
	}

	// Create service client
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	serviceClient, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		return "", err
	}

	// Get container client
	containerClient := serviceClient.NewContainerClient(containerName)

	// Upload blob
	blobClient := containerClient.NewBlockBlobClient(fileName)
	_, err = blobClient.Upload(context.Background(), file, nil)
	if err != nil {
		return "", err
	}

	return blobClient.URL(), nil
}
