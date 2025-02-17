package facades

import (
	"encoding/json"
	"go-metrics-alerting/internal/configs"

	"go-metrics-alerting/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetric_Success(t *testing.T) {
	// Создаем тестовый сервер, который будет имитировать ответ от настоящего API.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что метод запроса POST
		assert.Equal(t, "POST", r.Method)

		// Проверяем, что путь запроса /update/
		assert.Equal(t, "/update/", r.URL.Path)

		// Проверяем, что запрос содержит тело в формате JSON
		var requestBody types.MetricsRequest
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		assert.NoError(t, err)

		// Проверяем, что все поля соответствуют ожидаемым значениям
		assert.Equal(t, "cpu_usage", requestBody.ID) // ID
		assert.Equal(t, "gauge", requestBody.MType)  // MType
		assert.Nil(t, requestBody.Delta)             // Delta должно быть nil
		assert.Equal(t, 75.5, *requestBody.Value)    // Value должно быть 75.5

		// Отправляем успешный ответ
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	// Создаем конфигурацию с адресом нашего тестового сервера
	config := &configs.AgentConfig{
		Address: mockServer.URL,
	}

	// Создаем экземпляр MetricFacade с реальным HTTP клиентом
	client := resty.New()
	facade := NewMetricFacade(client, config)

	// Создаем метрику
	metric := types.MetricsRequest{
		ID:    "cpu_usage",
		MType: "gauge",
		Value: new(float64), // создаем указатель на значение
	}
	*metric.Value = 75.5 // Устанавливаем значение для Value

	// Пытаемся обновить метрику
	err := facade.UpdateMetric(metric)

	// Проверяем, что нет ошибок
	assert.NoError(t, err)
}

func TestUpdateMetric_ErrorStatusCode(t *testing.T) {
	// Создаем тестовый сервер, который будет имитировать ошибку (например, 500)
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Отправляем ошибку 500
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	// Создаем конфигурацию с адресом нашего тестового сервера
	config := &configs.AgentConfig{
		Address: mockServer.URL,
	}

	// Создаем экземпляр MetricFacade с реальным HTTP клиентом
	client := resty.New()
	facade := NewMetricFacade(client, config)

	// Создаем метрику
	metric := types.MetricsRequest{
		ID:    "cpu_usage",
		MType: "gauge",
		Value: new(float64),
	}
	*metric.Value = 75.5

	// Пытаемся обновить метрику
	err := facade.UpdateMetric(metric)

	// Проверяем, что произошла ошибка из-за неожидаемого кода статуса
	assert.Equal(t, ErrStatus, err)
}

func TestUpdateMetric_NetworkError(t *testing.T) {
	// Создаем конфигурацию с несуществующим адресом (для имитации ошибки сети)
	config := &configs.AgentConfig{
		Address: "http://localhost:9999", // Невалидный адрес
	}

	// Создаем экземпляр MetricFacade с реальным HTTP клиентом
	client := resty.New()
	facade := NewMetricFacade(client, config)

	// Создаем метрику
	metric := types.MetricsRequest{
		ID:    "cpu_usage",
		MType: "gauge",
		Value: new(float64),
	}
	*metric.Value = 75.5

	// Пытаемся обновить метрику
	err := facade.UpdateMetric(metric)

	// Проверяем, что произошла ошибка сети
	assert.Equal(t, ErrNetwork, err)
}
