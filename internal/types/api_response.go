package types

import "net/http"

// HTTPStatus - enum для статусов
type HTTPStatus int

const (
	StatusBadRequest HTTPStatus = http.StatusBadRequest
	StatusNotFound   HTTPStatus = http.StatusNotFound
	StatusOK         HTTPStatus = http.StatusOK
)

type APIResponse struct {
	Status HTTPStatus  `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type APIErrorResponse struct {
	Status  HTTPStatus `json:"status"`
	Message string     `json:"message"`
}

// newAPIError создает объект ошибки API
func NewAPIErrorResponse(status HTTPStatus, message string) *APIErrorResponse {
	return &APIErrorResponse{
		Status:  status,
		Message: message,
	}
}
