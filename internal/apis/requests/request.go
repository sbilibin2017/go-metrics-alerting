package requests

import (
	"net/http"

	"github.com/go-chi/chi"
)

func GetPathParam(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}
