package web

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// errorResponse represents the structure of an error message.
type errorResponse struct {
	Message string `json:"message"`
}

// newErrorResponse sends an error response with the specified status code and message.
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
