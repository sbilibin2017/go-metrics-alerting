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

// UpdateGaugeStrategy стратегия обновления значения для метрики типа Gauge
type UpdateGaugeStrategy struct {
	saver Saver
}

func (g *UpdateGaugeStrategy) Update(metric *domain.Metric) (*domain.Metric, bool) {
	ok := g.saver.Save(metric.ID, metric.Value)
	if !ok {
		return nil, false
	}
	return metric, true
}

// UpdateCounterStrategy стратегия обновления значения для метрики типа Counter
type UpdateCounterStrategy struct {
	saver  Saver
	getter Getter
}

func (c *UpdateCounterStrategy) Update(metric *domain.Metric) (*domain.Metric, bool) {
	currentValueStr, exists := c.getter.Get(metric.ID)
	if !exists {
		return nil, false
	}
	currentValue, ok := formatters.ParseInt64(currentValueStr)
	if !ok {
		return nil, false
	}
	newIntValue, ok := formatters.ParseInt64(metric.Value)
	if !ok {
		return nil, false
	}
	updatedValue := formatters.FormatInt64(currentValue + newIntValue)
	ok = c.saver.Save(metric.ID, updatedValue)
	if !ok {
		return nil, false
	}
	metric.Value = updatedValue
	return metric, true
}
