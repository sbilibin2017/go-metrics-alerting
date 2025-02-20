package facades

import (
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/domain"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// TestNewMetricFacade проверяет конструктор NewMetricFacade и добавление http:// к адресу
func TestNewMetricFacade(t *testing.T) {
	tests := []struct {
		name            string
		inputAddress    string
		expectedAddress string
	}{
		{"Address with http", "http://example.com", "http://example.com"},
		{"Address with https", "https://example.com", "https://example.com"},
		{"Address without prefix", "example.com", "http://example.com"},
		{"Address with port without prefix", "example.com:8080", "http://example.com:8080"},
		{"IP address without prefix", "192.168.1.1", "http://192.168.1.1"},
		{"IP address with port without prefix", "192.168.1.1:9000", "http://192.168.1.1:9000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем конфиг с тестируемым адресом
			config := &configs.AgentConfig{Address: tt.inputAddress}

			// Создаем экземпляр MetricFacade
			client := resty.New()
			logger, _ := zap.NewDevelopment()
			facade := NewMetricFacade(client, config, logger)

			// Проверяем, что адрес в конфиге соответствует ожидаемому
			assert.Equal(t, tt.expectedAddress, facade.config.Address, "Address was not correctly set")
		})
	}
}

// TestSuccessfulMetricSend проверяет успешную отправку метрики
func TestSuccessfulMetricSend(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	testLogger := zap.New(core)

	// Запускаем тестовый HTTP-сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := resty.New()
	config := &configs.AgentConfig{
		Address: server.URL,
	}

	facade := NewMetricFacade(client, config, testLogger)

	metrics := []*domain.Metric{
		{ID: "metric1", Value: "100"},
	}

	facade.UpdateMetrics(metrics)

	// Проверяем, что лог содержит сообщение об успешной отправке
	found := false
	for _, entry := range logs.All() {
		if entry.Message == "Metric sent successfully" {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected log message 'Metric sent successfully' not found")
}

// TestErrorLogging проверяет логирование ошибки при неудачной отправке метрики
func TestErrorLogging(t *testing.T) {
	core, logs := observer.New(zap.ErrorLevel)
	testLogger := zap.New(core)

	// Запускаем тестовый HTTP-сервер, который всегда возвращает ошибку
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := resty.New()
	config := &configs.AgentConfig{
		Address: server.URL,
	}

	facade := NewMetricFacade(client, config, testLogger)

	metrics := []*domain.Metric{
		{ID: "metric1", Value: "100"},
	}

	facade.UpdateMetrics(metrics)

	// Проверяем, что лог содержит сообщение об ошибке
	found := false
	for _, entry := range logs.All() {
		if entry.Message == "Error sending metric" {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected log message 'Error sending metric' not found")
}

// TestServerUnavailable проверяет, логируется ли ошибка, если сервер недоступен
func TestServerUnavailable(t *testing.T) {
	core, logs := observer.New(zap.ErrorLevel)
	testLogger := zap.New(core)

	// Указываем несуществующий адрес
	config := &configs.AgentConfig{
		Address: "http://127.0.0.1:9999",
	}

	client := resty.New()
	facade := NewMetricFacade(client, config, testLogger)

	metrics := []*domain.Metric{
		{ID: "metric1", Value: "100"},
	}

	facade.UpdateMetrics(metrics)

	// Проверяем, что лог содержит сообщение об ошибке
	found := false
	for _, entry := range logs.All() {
		if entry.Message == "Error sending metric" {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected log message 'Error sending metric' not found")
}
