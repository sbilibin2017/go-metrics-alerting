package responders

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

// RespondWithHTML рендерит HTML-шаблон и отправляет его клиенту, добавлено логирование.
func RespondWithHTML(c *gin.Context, statusCode int, tmplString string, data interface{}) {
	tmpl, err := template.New("response").Parse(tmplString)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.Header("Content-Type", "text/html")
	c.Writer.WriteHeader(statusCode)
	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err)
		return
	}
}
