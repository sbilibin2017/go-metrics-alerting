package types

// ApiClientResponse определяет структуру длф HTTP ответа.
type ApiClientResponse struct {
	StatusCode int
	Body       []byte
	Error      error
}
