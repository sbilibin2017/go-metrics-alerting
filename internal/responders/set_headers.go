package responders

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	textPlain string = "text/plain"
	textHTML  string = "text/html"
)

// SetHeaders is a helper function to set common response headers.
func setHeaders(c *gin.Context, contentType string) {
	c.Header("Content-Type", contentType+"; charset=utf-8")
	c.Header("Date", time.Now().UTC().Format(time.RFC1123))
}
