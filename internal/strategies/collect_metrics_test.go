package strategies

import (
	"fmt"
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaugeMetricsCollector_Collect(t *testing.T) {
	// Создаем экземпляр коллектора метрик типа Gauge
	collector := &GaugeMetricsCollector{}
	metrics := collector.Collect()

	// Проверяем, что метрики собраны
	assert.NotNil(t, metrics)
	assert.Equal(t, len(metrics), 14, "expected 14 metrics")

	// Проверяем, что метрика Alloc имеет правильный формат
	assert.Contains(t, metrics[0].ID, "Alloc")
	assert.Equal(t, metrics[0].MType, domain.Gauge)
	// Проверяем, что Value является строкой
	_, err := fmt.Sscanf(metrics[0].Value, "%f", new(float64))
	assert.NoError(t, err)

	// Проверяем другие метрики
	assert.Contains(t, metrics[1].ID, "BuckHashSys")
	assert.Equal(t, metrics[1].MType, domain.Gauge)

	// Проверка случайного значения
	assert.Contains(t, metrics[13].ID, "RandomValue")
	_, err = fmt.Sscanf(metrics[13].Value, "%f", new(float64))
	assert.NoError(t, err)
}

func TestCounterMetricsCollector_Collect(t *testing.T) {
	// Создаем экземпляр коллектора метрик типа Counter
	collector := &CounterMetricsCollector{}
	metrics := collector.Collect()

	// Проверяем, что метрика собрана
	assert.NotNil(t, metrics)
	assert.Equal(t, len(metrics), 1, "expected 1 metric")

	// Проверяем ID и значение метрики PollCount
	assert.Equal(t, metrics[0].ID, "PollCount")
	assert.Equal(t, metrics[0].MType, domain.Counter)

	// Проверяем, что значение метрики PollCount - строка с числовым значением
	var delta int64
	_, err := fmt.Sscanf(metrics[0].Value, "%d", &delta)
	assert.NoError(t, err)
	assert.Equal(t, delta, int64(1)) // Ожидаем, что это значение будет равно 1 в первый раз

	// Вызовем Collect еще раз и проверим, что значение увеличилось
	metrics = collector.Collect()
	assert.Equal(t, metrics[0].ID, "PollCount")
	var newDelta int64
	_, err = fmt.Sscanf(metrics[0].Value, "%d", &newDelta)
	assert.NoError(t, err)
	assert.Equal(t, newDelta, int64(2)) // Ожидаем, что значение увеличится до 2
}
