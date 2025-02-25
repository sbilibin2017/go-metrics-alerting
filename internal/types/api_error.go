package types

// APIError представляет ошибку, которая будет возвращена API.
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
