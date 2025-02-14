package responders

import (
	"go-metrics-alerting/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponder struct {
	C *gin.Context
}

func (e *ErrorResponder) Respond(err *types.APIErrorResponse) {
	setHeaders(e.C, textPlain)
	if err != nil {
		e.C.String(err.Code, err.Message)
	} else {
		e.C.String(http.StatusInternalServerError, "Internal Server Error")
	}
}
