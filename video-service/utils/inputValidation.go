package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
)

// ValidateFile checks if a file exists in the request and retrieves its details
func ValidateFile(c *gin.Context, fieldName string) (multipart.File, *multipart.FileHeader, error) {
	file, header, err := c.Request.FormFile(fieldName)
	if err != nil {
		return nil, nil, errors.New("failed to retrieve file: " + err.Error())
	}
	return file, header, nil
}

// ValidateRequiredField ensures that a required field is provided
func ValidateRequiredField(c *gin.Context, fieldName, errorMessage string) (string, error) {
	value := c.PostForm(fieldName)
	if value == "" {
		return "", errors.New(errorMessage)
	}
	return value, nil
}

// SaveTemporaryFile saves a file temporarily for processing
func SaveTemporaryFile(fileName string, file io.Reader) (string, error) {
	tmpFilePath := "/tmp/" + fileName
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		return "", fmt.Errorf("failed to write to temporary file: %w", err)
	}

	return tmpFilePath, nil
}
