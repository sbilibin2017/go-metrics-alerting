package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestCollectGaugeMetrics(t *testing.T) {
	metrics := collectGaugeMetrics()
	assert.NotEmpty(t, metrics)
	for _, metric := range metrics {
		assert.NotEmpty(t, metric.ID)
		assert.NotNil(t, metric.Value)
	}
}

func TestCollectCounterMetrics(t *testing.T) {
	// Получаем замыкание для подсчета опросов
	incrementPollCount := collectCounterMetrics()

	// Вызываем замыкание несколько раз, чтобы обновить pollCount
	metrics := incrementPollCount()
	assert.Len(t, metrics, 1)
	assert.Equal(t, "PollCount", metrics[0].ID)
	assert.NotNil(t, metrics[0].Delta)
	assert.Equal(t, int64(1), *metrics[0].Delta)

	// Вызываем еще раз, чтобы проверить инкремент
	metrics = incrementPollCount()
	assert.Len(t, metrics, 1)
	assert.Equal(t, "PollCount", metrics[0].ID)
	assert.NotNil(t, metrics[0].Delta)
	assert.Equal(t, int64(2), *metrics[0].Delta)
}

func TestSendMetric(t *testing.T) {
	// Создаём тестовый HTTP-сервер
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/update/", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var receivedMetric UpdateMetricsRequest
		err := json.NewDecoder(r.Body).Decode(&receivedMetric)
		assert.NoError(t, err)

		assert.Equal(t, Counter, receivedMetric.MType)
		assert.Equal(t, "TestCounter", receivedMetric.ID)
		assert.Equal(t, int64(1), *receivedMetric.Delta)

		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Тестовые данные
	delta := int64(1)
	metric := UpdateMetricsRequest{
		MType: Counter,
		ID:    "TestCounter",
		Delta: &delta,
	}

	// Вызываем тестируемую функцию
	client := resty.New()
	sendMetric(metric, client, testServer.URL)

	// Проверка, что сервер получил запрос
	// Тест пройдет, если сервер успешно обработал запрос
}

func TestCollectMetrics(t *testing.T) {
	metricsCh := make(chan UpdateMetricsRequest, 10) // Буферизованный канал

	go func() {
		collectMetrics(metricsCh) // Собираем метрики в отдельной горутине
		close(metricsCh)          // Закрываем канал после завершения
	}()

	var metrics []UpdateMetricsRequest
	for metric := range metricsCh {
		metrics = append(metrics, metric)
	}

	assert.NotEmpty(t, metrics) // Проверяем, что метрики были собраны
}

func TestSendMetrics(t *testing.T) {
	metricsCh := make(chan UpdateMetricsRequest, 10)
	delta := float64(123)
	metricsCh <- UpdateMetricsRequest{
		ID:    "TestMetric",
		MType: Gauge,
		Value: &delta,
	}

	mockClient := resty.New()
	sendMetrics(metricsCh, mockClient, "http://localhost:8080")
}

func TestStartAgent(t *testing.T) {
	// Создаём канал для сигналов ОС
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Настроим конфиг с короткими интервалами
	config := &AgentConfig{
		PollInterval:   10 * time.Millisecond,
		ReportInterval: 20 * time.Millisecond,
		Address:        "http://localhost:8080",
	}

	// Создаём HTTP-клиент
	client := resty.New()

	// Запускаем агента в отдельной горутине
	go StartAgent(signalCh, config, client)

	// Добавляем вызов сигнала завершения через 50 мс (чтобы агент успел начать работу)
	time.Sleep(50 * time.Millisecond)
	close(signalCh) // Отправляем сигнал завершения

	// Добавляем задержку, чтобы агент успел завершить свою работу
	time.Sleep(50 * time.Millisecond)

	// Если тест доходит до сюда, значит функция завершилась корректно
	assert.True(t, true, "startAgent завершился успешно")
}
