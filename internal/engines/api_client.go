package engines

import "github.com/go-resty/resty/v2"

// ApiClientInterface остается пустым, только как тип для интерфейса.
type ApiClientInterface interface {
	R() *resty.Request
}

// Реализация интерфейса ApiClientInterface для реального клиента.
type ApiClient struct {
	client *resty.Client
}

// Конструктор для создания реального клиента.
func NewApiClient() *ApiClient {
	return &ApiClient{
		client: resty.New(),
	}
}

// Доступ к клиенту Resty с помощью метода R()
func (a *ApiClient) R() *resty.Request {
	return a.client.R()
}
