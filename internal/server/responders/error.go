package responders

import (
	"github.com/gin-gonic/gin"
)

// RespondWithError отправляет JSON-ответ с ошибкой и логирует ее.
func RespondWithError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}
