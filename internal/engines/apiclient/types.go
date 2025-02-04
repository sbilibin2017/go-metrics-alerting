package apiclient

import (
	"time"
)

// Response содержит HTTP-ответ.
type ApiResponse struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

// ClientOptions задает параметры для API-клиента.
type ApiClientOptions struct {
	Timeout    time.Duration
	RetryCount int
}
