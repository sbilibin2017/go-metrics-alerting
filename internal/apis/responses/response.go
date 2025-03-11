package responses

import (
	"encoding/json"
	"net/http"
)

func NotFoundResponse(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusNotFound)
}

func BadRequestResponse(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func InternalServerErrorResponse(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func sendErrorResponse(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}

func TextResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func JsonResponse(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		InternalServerErrorResponse(w, err)
	}
}
