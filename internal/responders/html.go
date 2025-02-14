package responders

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

type HTMLHandler struct {
	C *gin.Context
}

// Respond принимает произвольный HTML-шаблон и данные для рендеринга.
func (h *HTMLHandler) Respond(templateContent string, data interface{}) {
	setHeaders(h.C, textHTML)
	tmpl, err := template.New("dynamicTemplate").Parse(templateContent)
	if err != nil {
		h.C.String(http.StatusInternalServerError, "Template parsing failed: "+err.Error())
		return
	}
	h.C.Writer.WriteHeader(http.StatusOK)
	tmpl.Execute(h.C.Writer, data)
}
