package types

// APIError представляет ошибку, которую API возвращает в ответе
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
