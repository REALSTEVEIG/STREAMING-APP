package utils

import (
	"github.com/gin-gonic/gin"
)

// RespondWithError sends an error response in a consistent format
func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

// RespondWithSuccess sends a success response in a consistent format
func RespondWithSuccess(c *gin.Context, statusCode int, data gin.H) {
	c.JSON(statusCode, data)
}
