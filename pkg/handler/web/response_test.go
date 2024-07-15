package web

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestNewErrorResponse(t *testing.T) {
	var logBuffer bytes.Buffer
	logrus.SetOutput(&logBuffer)
	defer func() {
		logrus.SetOutput(nil)
	}()

	router := gin.Default()

	// Define a test route that triggers the newErrorResponse function
	router.GET("/test", func(c *gin.Context) {
		newErrorResponse(c, 500, "Internal Server Error")
	})

	// test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)

	var resp errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Internal Server Error", resp.Message)

	// Check the log output
	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "Internal Server Error")
}
