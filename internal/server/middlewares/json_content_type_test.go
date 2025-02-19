package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJSONContentTypeMiddleware_ValidContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(JSONContentTypeMiddleware())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	reqBody := []byte(`{"key":"value"}`)
	req, _ := http.NewRequest("POST", "/test", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"message":"success"`)
}

func TestJSONContentTypeMiddleware_InvalidContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(JSONContentTypeMiddleware())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	reqBody := []byte(`{"key":"value"}`)
	req, _ := http.NewRequest("POST", "/test", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "text/plain")
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnsupportedMediaType, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"error":"Unsupported Media Type"`)
}
