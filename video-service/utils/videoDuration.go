package utils

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)


// IsVideoContentType checks if the content type represents a video
func IsVideoContentType(contentType string) bool {
	videoContentTypes := []string{"video/mp4", "video/avi", "video/mpeg", "video/quicktime", "video/x-matroska"}
	for _, v := range videoContentTypes {
		if contentType == v {
			return true
		}
	}
	return false
}

// CalculateVideoDuration calculates the duration of a video file in seconds
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
