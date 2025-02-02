package utils

import (
	"net/http"
	"twittir-go/internal/dto"

	"github.com/gin-gonic/gin"
)

// Helper untuk merespon error
func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, dto.ErrorResponse{Status: http.StatusText(statusCode), Message: message})
}

func RespondWithSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, dto.SuccessResponse{Status: http.StatusText(statusCode), Data: data})
}
