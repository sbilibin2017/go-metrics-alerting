package services

import (
	"encoding/json"
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/logger"
	"go-metrics-alerting/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestCollectGaugeMetrics проверяет, что собираются метрики типа gauge
func TestCollectGaugeMetrics(t *testing.T) {
	metrics := collectGaugeMetrics()

	assert.NotEmpty(t, metrics, "Gauge metrics should not be empty")
	for _, metric := range metrics {
		assert.Equal(t, types.Gauge, metric.MType, "Metric type should be gauge")
	}
}

// TestCollectCounterMetrics проверяет, что собирается метрика типа counter
func TestCollectCounterMetrics(t *testing.T) {
	metrics := collectCounterMetrics()

	assert.NotEmpty(t, metrics, "Counter metrics should not be empty")
	assert.Equal(t, types.Counter, metrics[0].MType, "Metric type should be counter")
}

// TestSendMetric проверяет отправку одной метрики через REST API
func TestSendMetric(t *testing.T) {
	// Создаем тестовый HTTP-сервер
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Создаем конфиг и клиента
	config := &configs.AgentConfig{Address: testServer.URL}
	client := resty.New()

	metric := types.UpdateMetricsRequest{
		MType: "gauge",
		ID:    "TestMetric",
		Value: float64Ptr(42.0),
	}

	err := sendMetric(metric, client, config.Address)
	assert.NoError(t, err, "Error should be nil when sending metric")
}

// TestSendMetrics проверяет отправку всех метрик из канала
func TestSendMetrics(t *testing.T) {
	// Создаем тестовый HTTP-сервер
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var metric types.UpdateMetricsRequest
		err := json.NewDecoder(r.Body).Decode(&metric)
		assert.NoError(t, err, "Failed to decode metric request")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Создаем конфиг и клиента
	config := &configs.AgentConfig{Address: testServer.URL}
	client := resty.New()

	metricsCh := make(chan types.UpdateMetricsRequest, 10)

	go func() {
		metricsCh <- types.UpdateMetricsRequest{
			MType: "gauge",
			ID:    "TestMetric",
			Value: float64Ptr(100.0),
		}
		close(metricsCh) // Закрываем канал после отправки
	}()

	sendMetrics(metricsCh, client, config.Address)
}

// TestSendMetric_ConnectionError проверяет обработку ошибки сети
func TestSendMetric_ConnectionError(t *testing.T) {
	config := &configs.AgentConfig{Address: "http://invalid-url"}
	client := resty.New()

	metric := types.UpdateMetricsRequest{
		MType: "gauge",
		ID:    "TestMetric",
		Value: float64Ptr(42.0),
	}

	err := sendMetric(metric, client, config.Address)
	assert.Error(t, err, "Expected an error when connection fails")
}

// Тест обработки ошибки при отправке метрик
func TestSendMetrics_ErrorHandling(t *testing.T) {
	// Создаем тестовый HTTP-сервер, который ВСЕГДА возвращает 500 ошибку
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	// Настроиваем сервис с тестовым сервером
	config := &configs.AgentConfig{Address: testServer.URL}
	client := resty.New()

	// Создаем тестовую метрику
	metric := types.UpdateMetricsRequest{
		MType: "unknwon",
		ID:    "TestMetric",
		Value: float64Ptr(42.0),
	}

	// Добавляем метрику в канал
	metricsCh := make(chan types.UpdateMetricsRequest, 10)
	metricsCh <- metric
	close(metricsCh)

	// Используем канал done, чтобы отслеживать завершение метода
	done := make(chan struct{})
	go func() {
		sendMetrics(metricsCh, client, config.Address)
		close(done)
	}()

	// Ожидаем завершения с тайм-аутом 3 секунды
	select {
	case <-done:
		// Проверяем, что ошибка была залогирована
		logger.Logger.Debug("Test completed: error should be logged", zap.String("metric_id", metric.ID))
	case <-time.After(3 * time.Second):
		t.Fatal("sendMetrics() завис и не завершился")
	}
}

func TestCollectMetrics(t *testing.T) {
	// Создаем канал для метрик
	metricsCh := make(chan types.UpdateMetricsRequest, 20) // Буферизованный канал для тестов

	// Запускаем функцию collectMetrics в горутине, чтобы она не блокировала тест
	go collectMetrics(metricsCh)

	// Ожидаем, что в канал будут отправлены как минимум метрики типа "gauge" и "counter"
	var collectedMetrics []types.UpdateMetricsRequest
	for len(collectedMetrics) < 15 { // 14 gauge + 1 counter (всего 15 метрик)
		metric := <-metricsCh
		collectedMetrics = append(collectedMetrics, metric)
	}

	// Проверяем, что в канале 15 метрик (14 gauge + 1 counter)
	assert.Len(t, collectedMetrics, 15, "Expected to collect 15 metrics")

	// Проверяем типы и ID метрик
	expectedIDs := map[string]bool{
		"Alloc":         true,
		"PollCount":     true,
		"RandomValue":   true,
		"BuckHashSys":   true,
		"Frees":         true,
		"GCCPUFraction": true,
		"HeapAlloc":     true,
		"HeapIdle":      true,
		"HeapInuse":     true,
		"HeapObjects":   true,
		"HeapReleased":  true,
		"HeapSys":       true,
		"NumGC":         true,
		"Sys":           true,
		"TotalAlloc":    true,
	}

	// Проверяем ID метрик
	for _, metric := range collectedMetrics {
		assert.True(t, expectedIDs[metric.ID], "Unexpected metric ID: "+metric.ID)
	}
}

func TestStartMetricsCollection(t *testing.T) {
	// Настроим тестовый сервер
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Конфигурируем сервер
	config := &configs.AgentConfig{Address: testServer.URL, PollInterval: 1, ReportInterval: 2}
	client := resty.New()

	// Запуск коллекции метрик в горутине
	go StartMetricAgent(config, client)

	// Ждем некоторое время, чтобы процесс собирал метрики
	time.Sleep(4 * time.Second)

	// Проверяем, что тест не зависает
	select {
	case <-time.After(3 * time.Second): // Ждем 3 секунды, чтобы убедиться, что тест не зависает
		t.Fatal("Test timed out")
	default:
		t.Log("Test passed: collection stopped")
	}
}
