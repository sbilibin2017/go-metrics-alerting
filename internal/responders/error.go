package responders

import (
	"github.com/gin-gonic/gin"
)

// respondWithError отправляет JSON-ответ с ошибкой.
func RespondWithError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}
