package strategies

import (
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/formatters"
)

// MetricSaver интерфейс для сохранения метрик.
type Saver interface {
	Save(id string, value string) bool
}

// MetricGetter интерфейс для получения метрик.
type Getter interface {
	Get(id string) (string, bool)
}

// KeyEncoder интерфейс для кодирования ключей
type KeyEncoder interface {
	Encode(id string, mtype string) string
}

// UpdateGaugeStrategy стратегия обновления значения для метрики типа Gauge
type UpdateGaugeStrategy struct {
	saver   Saver
	encoder KeyEncoder
}

// NewUpdateGaugeStrategy создает новый объект стратегии обновления метрики типа Gauge
func NewUpdateGaugeStrategy(saver Saver, encoder KeyEncoder) *UpdateGaugeStrategy {
	return &UpdateGaugeStrategy{
		saver:   saver,
		encoder: encoder,
	}
}

func (g *UpdateGaugeStrategy) Update(metric *domain.Metric) (*domain.Metric, bool) {
	// Генерируем ключ с помощью KeyEncoder
	key := g.encoder.Encode(metric.ID, string(metric.MType))

	// Сохраняем значение метрики по сгенерированному ключу
	ok := g.saver.Save(key, metric.Value)
	if !ok {
		return nil, false
	}
	return metric, true
}

// UpdateCounterStrategy стратегия обновления значения для метрики типа Counter
type UpdateCounterStrategy struct {
	saver   Saver
	getter  Getter
	encoder KeyEncoder
}

// NewUpdateCounterStrategy создает новый объект стратегии обновления метрики типа Counter
func NewUpdateCounterStrategy(saver Saver, getter Getter, encoder KeyEncoder) *UpdateCounterStrategy {
	return &UpdateCounterStrategy{
		saver:   saver,
		getter:  getter,
		encoder: encoder,
	}
}

func (c *UpdateCounterStrategy) Update(metric *domain.Metric) (*domain.Metric, bool) {
	// Генерируем ключ с помощью KeyEncoder
	key := c.encoder.Encode(metric.ID, string(metric.MType))

	// Получаем текущее значение метрики по ключу
	currentValueStr, exists := c.getter.Get(key)
	if !exists {
		return nil, false
	}
	currentValue, ok := formatters.ParseInt64(currentValueStr)
	if !ok {
		return nil, false
	}

	// Парсим новое значение метрики и обновляем
	newIntValue, ok := formatters.ParseInt64(metric.Value)
	if !ok {
		return nil, false
	}
	updatedValue := formatters.FormatInt64(currentValue + newIntValue)

	// Сохраняем обновленное значение по ключу
	ok = c.saver.Save(key, updatedValue)
	if !ok {
		return nil, false
	}

	// Обновляем метрику
	metric.Value = updatedValue
	return metric, true
}
