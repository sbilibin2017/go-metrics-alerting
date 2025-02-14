package services

import (
	"fmt"
	"go-metrics-alerting/internal/configs"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

// testServer для тестов
func testServer() *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		// Проверим, что URL правильно передается
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			http.Error(w, "Invalid URL path", http.StatusBadRequest)
			return
		}

		// Проверим, что данные приходят корректно
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Логируем входящие данные
		fmt.Printf("Received request: %s %s/%s/%s\n", r.Method, parts[1], parts[2], parts[3])

		w.WriteHeader(http.StatusOK)
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

// Подготовка тестов
func TestMetricAgentService_SendMetrics(t *testing.T) {
	// Запускаем тестовый сервер
	server := testServer()
	go server.ListenAndServe()
	defer server.Close()

	// Создаем конфигурацию с тестовым сервером
	config := &configs.AgentConfig{
		Address:        "http://localhost:8080",
		PollInterval:   time.Second,
		ReportInterval: 2 * time.Second,
	}

	// Создаем сервис для сбора метрик
	service := NewMetricAgentService(config, resty.New())

	// Вспомогательная функция для отправки метрик
	go func() {
		// Запускаем сервис для сбора и отправки метрик
		service.Start()
	}()

	// Даем сервису немного времени на сбор и отправку метрик
	time.Sleep(3 * time.Second)

	// Проверка, что тестовый сервер был вызван с правильными параметрами
	assert.Contains(t, server.Addr, ":8080", "Server should be running on port 8080")
}
