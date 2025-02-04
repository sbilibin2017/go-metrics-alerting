package apiclient

import (
	"context"
	"errors"
	"net"

	"github.com/go-resty/resty/v2"
)

// ApiClientEngine реализует интерфейс HTTP-клиента.
type ApiClientEngine struct {
	resty *resty.Client
}

// NewClient создает новый API-клиент.
func NewClient(options ApiClientOptions) *ApiClientEngine {
	client := resty.New().
		SetTimeout(options.Timeout).
		SetRetryCount(options.RetryCount)

	return &ApiClientEngine{resty: client}
}

// Get выполняет GET-запрос.
func (c *ApiClientEngine) Get(path string, query map[string]string, headers map[string]string) (*ApiResponse, error) {
	resp, err := c.resty.R().
		SetQueryParams(query).
		SetHeaders(headers).
		Get(path)
	return c.classifyError(resp, err)
}

// Post выполняет POST-запрос.
func (c *ApiClientEngine) Post(path string, query map[string]string, body any, headers map[string]string) (*ApiResponse, error) {
	resp, err := c.resty.R().
		SetQueryParams(query).
		SetHeaders(headers).
		SetBody(body).
		Post(path)
	return c.classifyError(resp, err)
}

// classifyError обрабатывает ошибки и возвращает соответствующие значения.
func (c *ApiClientEngine) classifyError(resp *resty.Response, err error) (*ApiResponse, error) {
	if isTimeoutError(err) {
		return nil, ErrTimeout
	}
	if err != nil {
		return nil, ErrRequestFailed
	}
	if resp == nil || resp.StatusCode() < 200 || resp.StatusCode() >= 300 || len(resp.Body()) == 0 {
		return nil, ErrInvalidResponse
	}
	return &ApiResponse{
		StatusCode: resp.StatusCode(),
		Headers:    resp.Header(),
		Body:       resp.Body(),
	}, nil
}

// isTimeoutError проверяет, является ли ошибка тайм-аутом.
func isTimeoutError(err error) bool {
	var netErr net.Error
	return errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &netErr) && netErr.Timeout())
}
