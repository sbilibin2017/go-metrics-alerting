package responders

import "github.com/gin-gonic/gin"

// ErrorResponderInterface определяет метод для обработки ошибок
type ResponderInterface interface {
	Respond(c *gin.Context, statusCode int, message string)
}
