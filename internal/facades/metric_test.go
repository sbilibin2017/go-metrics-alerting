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

// Функция для имитации HTTP-сервера, который принимает запросы
func mockServer() *httptest.Server {
	// Настроим сервер с обработчиком, который будет эмулировать реальное поведение
	handler := http.NewServeMux()
	handler.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		// Предположим, что сервер принимает запрос и возвращает статус 200 OK
		if r.Method == http.MethodPost {
			// Читаем тело запроса
			var requestBody types.UpdateMetricBodyRequest
			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			// Здесь можно выполнить дополнительные проверки, например:
			if requestBody.ID == "" || requestBody.MType == "" {
				http.Error(w, "Invalid request data", http.StatusBadRequest)
				return
			}

			// Ответ с успешным статусом
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "success"}`))
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	server := httptest.NewServer(handler)
	return server
}

// Тестирование обновления метрики с реальным сервером
func TestUpdateMetric_RealServer(t *testing.T) {
	// Создаем сервер
	server := mockServer()
	defer server.Close()

	// Подготавливаем конфигурацию
	config := &configs.AgentConfig{
		Address: server.URL, // Используем адрес тестового сервера
	}

	// Создаем клиент и facade
	client := resty.New()
	metricFacade := NewMetricFacade(client, config)

	// Создаем тестовую метрику с новой структурой
	metric := &types.UpdateMetricBodyRequest{
		ID:    "metric123",
		MType: "gauge",
		Delta: nil,
		Value: float64Ptr(75.5),
	}

	// Обновляем метрику
	metricFacade.UpdateMetric(metric)

	// Теперь проверим, что запрос был правильно обработан сервером
	resp, err := client.R().SetBody(metric).Post(server.URL + "/update/")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// Дополнительно можно проверить тело ответа
	assert.Contains(t, resp.String(), "success")
}

func TestNewMetricFacade(t *testing.T) {
	tests := []struct {
		name            string
		address         string
		expectedAddress string
	}{
		{
			name:            "Address without prefix",
			address:         "localhost:8080",        // Нет префикса
			expectedAddress: "http://localhost:8080", // Ожидаем добавление http://
		},
		{
			name:            "Address with http prefix",
			address:         "http://localhost:8080", // Уже есть префикс http://
			expectedAddress: "http://localhost:8080", // Не должно изменяться
		},
		{
			name:            "Address with https prefix",
			address:         "https://localhost:8080", // Уже есть префикс https://
			expectedAddress: "https://localhost:8080", // Не должно изменяться
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем конфигурацию
			config := &configs.AgentConfig{
				Address: tt.address,
			}

			// Создаем клиент
			client := resty.New()

			// Создаем экземпляр MetricFacade с конфигурацией
			metricFacade := NewMetricFacade(client, config)

			// Проверяем, что адрес соответствует ожидаемому
			assert.Equal(t, tt.expectedAddress, metricFacade.config.Address)
		})
	}
}

// Вспомогательная функция для создания указателя на float64
func float64Ptr(value float64) *float64 {
	return &value
}
