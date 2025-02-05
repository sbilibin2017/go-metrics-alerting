package apiclient

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Response представляет ответ на POST-запрос
type ApiResponse struct {
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
func (r *restyClient) Post(urlString string, headers map[string]string) (*ApiResponse, error) {
	// Добавляем схему, если она отсутствует
	parsedURL, err := addSchemeIfMissing(urlString)
	if err != nil {
		return nil, err
	}

	// Создаем новый запрос
	req := r.client.R()

	// Устанавливаем заголовки
	for key, value := range headers {
		req.SetHeader(key, value)
	}

	// Выполняем POST-запрос
	resp, err := req.Post(parsedURL)
	if err != nil {
		return nil, err
	}

	// Возвращаем ответ
	return &ApiResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.String(),
	}, nil
}

// Функция, которая добавляет схему, если она отсутствует
func addSchemeIfMissing(urlString string) (string, error) {
	// Если в URL нет схемы, добавляем http:// по умолчанию
	if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
		urlString = "http://" + urlString
	}

	// Теперь парсим URL
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	return parsedURL.String(), nil
}
