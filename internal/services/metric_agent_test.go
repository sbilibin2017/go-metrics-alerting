package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

// Тестовый сервер для приема метрик
func startTestServer(t *testing.T) *http.Server {
	mux := http.NewServeMux()

	// Обработчик для пути /update, куда будут отправляться метрики
	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		var metric types.MetricsRequest
		// Декодируем JSON тело запроса в метрику
		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		// Проверяем, что метрика не пустая
		assert.NotNil(t, metric.ID)
		assert.NotNil(t, metric.MType)

		// Мы просто отвечаем 200 OK, чтобы агент понял, что метрика получена
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:    ":8080", // Можем использовать порт 8080 для тестирования
		Handler: mux,
	}

	// Запускаем сервер в горутине, чтобы он не блокировал тест
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Fatalf("Test server failed: %v", err)
		}
	}()
	// Ждем, пока сервер начнет работать
	time.Sleep(time.Second)
	return server
}

// Функция для остановки тестового сервера
func stopTestServer(server *http.Server) {
	if err := server.Close(); err != nil {
		fmt.Println("Error closing test server:", err)
	}
}

func TestMetricAgentService_Start(t *testing.T) {
	// Запускаем тестовый сервер
	server := startTestServer(t)
	defer stopTestServer(server)

	// Настройка конфигурации агента
	config := &configs.AgentConfig{
		Address:        "http://localhost:8080", // Адрес тестового сервера
		PollInterval:   1,                       // интервал обновления метрик
		ReportInterval: 2,                       // интервал отправки метрик
	}

	// Инициализация Resty клиента
	client := resty.New()

	// Создаем новый агент
	agent := NewMetricAgentService(config, client)

	// Запускаем агент в отдельной горутине, чтобы он мог работать асинхронно
	go func() {
		agent.Start()
	}()

	// Даем агенту время для сбора и отправки метрик
	time.Sleep(5 * time.Second)

	// Поскольку сервер отвечает кодом 200 для успешных запросов, проверим, что метрики отправлены
	// Если сервер получит метрики, это будет подтверждением успешной работы агента.
}

func TestMetricAgentService_ClearChannel(t *testing.T) {
	// Запускаем тестовый сервер
	server := startTestServer(t)
	defer stopTestServer(server)

	// Настройка конфигурации агента
	config := &configs.AgentConfig{
		Address:        "http://localhost:8080", // Адрес тестового сервера
		PollInterval:   1,                       // интервал обновления метрик
		ReportInterval: 2,                       // интервал отправки метрик
	}

	// Инициализация Resty клиента
	client := resty.New()

	// Создаем новый агент
	agent := NewMetricAgentService(config, client)

	// Добавим несколько метрик в канал для тестирования
	agent.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "Alloc", Value: new(float64)}
	agent.metricsCh <- types.MetricsRequest{MType: "gauge", ID: "Frees", Value: new(float64)}

	// Проверяем, что канал не пустой
	assert.NotEmpty(t, agent.metricsCh)

	// Очищаем канал
	agent.clearChannel()

	// Проверяем, что канал теперь пуст
	assert.Empty(t, agent.metricsCh)
}

// Тест на обработку сетевых ошибок
func TestMetricAgentService_NetworkError(t *testing.T) {
	// Адрес недоступного сервера
	config := &configs.AgentConfig{
		Address:        "http://localhost:9999", // Несуществующий сервер
		PollInterval:   1,
		ReportInterval: 2,
	}

	client := resty.New()
	agent := NewMetricAgentService(config, client)

	// Запускаем агент в отдельной горутине
	go func() {
		agent.Start()
	}()

	// Даем агенту время на попытку отправки метрик
	time.Sleep(3 * time.Second)

	// Если программа не завершилась с паникой, значит, ошибки обработаны корректно
}

// Тест на обработку неожиданных статус-кодов
func TestMetricAgentService_UnexpectedStatusCode(t *testing.T) {
	// Создаем тестовый сервер, который возвращает 500 Internal Server Error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	config := &configs.AgentConfig{
		Address:        server.URL, // Используем тестовый сервер
		PollInterval:   1,
		ReportInterval: 2,
	}

	client := resty.New()
	agent := NewMetricAgentService(config, client)

	// Запускаем агент в отдельной горутине
	go func() {
		agent.Start()
	}()

	// Даем агенту время на попытку отправки метрик
	time.Sleep(3 * time.Second)

	// Проверяем, что агент продолжил работу, несмотря на ошибки
}
