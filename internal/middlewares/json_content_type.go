package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	JsonContentType         string = "application/json"
	ErrUnsupportedMediaType string = "Unsupported Media Type"
)

// JSONContentTypeMiddleware проверяет, что Content-Type = application/json для POST-запросов
func JSONContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			if c.ContentType() != JsonContentType {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": ErrUnsupportedMediaType})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
