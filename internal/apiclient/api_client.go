package apiclient

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// APIResponse представляет ответ на POST-запрос
type APIResponse struct {
	StatusCode int
	Body       string
}

// RestyClient структура, которая реализует интерфейс PostRequester
type restyClient struct {
	client *resty.Client
}

// NewRestyClient создает новый экземпляр RestyClient
func NewRestyClient() *restyClient {
	return &restyClient{
		client: resty.New(),
	}
}

// Post выполняет POST-запрос, используя Resty
// Post выполняет POST-запрос, используя Resty
func (r *restyClient) Post(urlString string, headers map[string]string) (*APIResponse, error) {
	// Если в URL нет схемы, добавляем http:// по умолчанию
	if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
		urlString = "http://" + urlString
	}

	// Создаем новый запрос
	req := r.client.R()

	// Если в заголовках нет Content-Type, устанавливаем его по умолчанию
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json"
	}

	// Устанавливаем заголовки
	for key, value := range headers {
		req.SetHeader(key, value)
	}

	// Выполняем POST-запрос
	resp, err := req.Post(urlString)
	if err != nil {
		return nil, fmt.Errorf("error while making POST request: %v", err)
	}

	// Проверяем статус код ответа
	if resp.StatusCode() >= 400 {
		return &APIResponse{
			StatusCode: resp.StatusCode(),
			Body:       resp.String(),
		}, fmt.Errorf("request failed with status %d", resp.StatusCode())
	}

	// Возвращаем успешный ответ
	return &APIResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.String(),
	}, nil
}
