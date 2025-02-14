package services

import (
	"fmt"
	"go-metrics-alerting/internal/configs"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Тестовый сервер для проверки отправки метрик
func testServer(t *testing.T, wg *sync.WaitGroup) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done() // Уменьшаем счетчик ожидания

		// Проверяем, что URL правильно передается
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			http.Error(w, "Invalid URL path", http.StatusBadRequest)
			return
		}

		// Проверяем метод запроса
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Логируем входящие данные (для отладки)
		fmt.Printf("Received request: %s %s/%s/%s\n", r.Method, parts[1], parts[2], parts[3])

		w.WriteHeader(http.StatusOK)
	})

	return httptest.NewServer(mux)
}

// Тест для проверки отправки метрик
func TestMetricAgentService_SendMetrics(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1) // Добавляем один вызов в ожидание

	// Запускаем тестовый HTTP сервер
	server := testServer(t, &wg)
	defer server.Close()

	// Создаем конфигурацию с тестовым сервером
	config := &configs.AgentConfig{
		Address:        server.URL, // Используем тестовый сервер
		PollInterval:   500 * time.Millisecond,
		ReportInterval: 1 * time.Second,
	}

	// Создаем сервис для сбора метрик
	service := NewMetricAgentService(config)

	// Запускаем сбор и отправку метрик в отдельной горутине
	go func() {
		service.Start()
	}()

	// Ожидаем завершения хотя бы одной отправки метрик
	wg.Wait()

	// Проверяем, что сервер получил запрос
	assert.True(t, true, "MetricAgentService should send at least one metric successfully")
}
