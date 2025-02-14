package responders

import (
	"github.com/gin-gonic/gin"
)

type SuccessResponder struct {
	C *gin.Context
}

func (s *SuccessResponder) Respond(message string, statusCode int) {
	setHeaders(s.C, textPlain)
	s.C.String(statusCode, message)
}
