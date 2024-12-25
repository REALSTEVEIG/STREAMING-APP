package utils

import (
	"context"
	"fmt"
	"io"
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// UploadFileToS3 uploads a file to an AWS S3 bucket
func UploadFileToS3(uploader *manager.Uploader, bucket, fileName, contentType string, file io.Reader, acl types.ObjectCannedACL) (string, error) {
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String(contentType),
		ACL:         acl,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}
	return result.Location, nil
}

func IsVideoContentType(contentType string) bool {
	videoContentTypes := []string{"video/mp4", "video/avi", "video/mpeg", "video/quicktime", "video/x-matroska"}
	for _, v := range videoContentTypes {
		if contentType == v {
			return true
		}
	}
	return false
}

// IsImageContentType checks if the content type represents an Image
func IsImageContentType(contentType string) bool {
	imageContentTypes := []string{"image/jpeg", "image/png", "image/webp"}
	for _, v := range imageContentTypes {
		if contentType == v {
			return true
		}
	}
	return false
}

// CalculateVideoDuration ca`lculates the duration of a video file in seconds
func CalculateVideoDuration(filePath string) (int, error) {
	// Run ffprobe to get video metadata
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "json", filePath)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to run ffprobe: %w", err)
	}

	// Parse the output JSON
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return 0, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	// Extract duration
	format, ok := result["format"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid ffprobe output format")
	}
	durationStr, ok := format["duration"].(string)
	if !ok {
		return 0, fmt.Errorf("duration not found in ffprobe output")
	}

	// Convert duration to an integer
	duration, err := strconv.ParseFloat(strings.TrimSpace(durationStr), 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return int(duration), nil
}

