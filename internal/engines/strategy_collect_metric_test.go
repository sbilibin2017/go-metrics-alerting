package engines

import (
	"go-metrics-alerting/internal/types"
	"runtime"
	"testing"
)

// Мок для MemoryStatsProvider
type MockMemoryStatsProvider struct{}

func (m *MockMemoryStatsProvider) ReadMemStats(ms *runtime.MemStats) {
	ms.Alloc = 123456
	ms.BuckHashSys = 654321
	ms.Frees = 1000
	// Добавьте другие поля по необходимости
}

func TestGaugeCollectionStrategyEngine_CollectMetrics(t *testing.T) {
	// Используем мок вместо реального считывания статистики
	engine := &GaugeCollectionStrategyEngine{
		memStatsProvider: &MockMemoryStatsProvider{},
	}

	metrics := engine.CollectMetrics()

	// Check if we have 28 metrics for the Gauge type (including "RandomValue")
	if len(metrics) != 28 {
		t.Fatalf("Expected 28 metrics, but got %d", len(metrics))
	}

	// Check if the "Alloc" metric is correctly generated
	expectedMetric := map[string]interface{}{
		"Type":  string(types.GaugeType),
		"Name":  "Alloc",
		"Value": FormatFloat(float64(123456)), // Мокированное значение
	}

	// Проверка для первой метрики
	if metrics[0].(map[string]interface{})["Name"] != expectedMetric["Name"] {
		t.Errorf("Expected metric name 'Alloc', but got '%v'", metrics[0].(map[string]interface{})["Name"])
	}
}

func TestCounterCollectionStrategyEngine_CollectMetrics(t *testing.T) {
	engine := &CounterCollectionStrategyEngine{}

	// Первый вызов: PollCount должен быть 0
	metrics := engine.CollectMetrics()
	if len(metrics) != 1 {
		t.Fatalf("Expected 1 metric, but got %d", len(metrics))
	}

	expectedMetric := map[string]interface{}{
		"Type":  string(types.CounterType),
		"Name":  "PollCount",
		"Value": FormatInt(0),
	}

	// Проверка значения на первом вызове
	if metrics[0].(map[string]interface{})["Value"] != expectedMetric["Value"] {
		t.Errorf("Expected PollCount value '%v', but got '%v'", expectedMetric["Value"], metrics[0].(map[string]interface{})["Value"])
	}

	// Второй вызов: PollCount должен быть 1
	metrics = engine.CollectMetrics()
	expectedMetric["Value"] = FormatInt(1)

	// Проверка значения на втором вызове
	if metrics[0].(map[string]interface{})["Value"] != expectedMetric["Value"] {
		t.Errorf("Expected PollCount value '%v', but got '%v'", expectedMetric["Value"], metrics[0].(map[string]interface{})["Value"])
	}
}

func TestMetricCollectionWithRandomValue(t *testing.T) {
	engine := &GaugeCollectionStrategyEngine{
		memStatsProvider: &MockMemoryStatsProvider{},
	}
	metrics := engine.CollectMetrics()

	// Проверяем, что среди метрик есть RandomValue
	var randomValue float64
	for _, metric := range metrics {
		metricMap := metric.(map[string]interface{})
		if metricMap["Name"] == "RandomValue" {
			// Проверяем, что RandomValue является числом типа float64
			if v, ok := metricMap["Value"].(float64); ok {
				randomValue = v
			} else {
				t.Errorf("Expected RandomValue to be of type float64, but got %T", metricMap["Value"])
			}
		}
	}

	// Дополнительная проверка, что randomValue не равно нулю
	if randomValue == 0 {
		t.Errorf("RandomValue should not be 0")
	}

	// Дополнительная проверка, что randomValue лежит в пределах ожидаемых значений (0 <= randomValue < 100)
	if randomValue < 0 || randomValue >= 100 {
		t.Errorf("RandomValue should be between 0 and 100, but got %v", randomValue)
	}
}
