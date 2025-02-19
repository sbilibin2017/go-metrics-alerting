package responders

import (
	"github.com/gin-gonic/gin"
)

// respondWithSuccess отправляет JSON-ответ с успешными данными и логирует информацию о запросе.
func RespondWithSuccess(c *gin.Context, statusCode int, payload interface{}) {
	c.JSON(statusCode, payload)
}
