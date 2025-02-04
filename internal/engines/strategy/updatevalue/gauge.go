package updatevalue

import (
	"go-metrics-alerting/internal/engines/numberprocessor"
)

// GaugeUpdateStrategyEngine handles gauge metrics.
type UpdateGaugeValueStrategyEngine[T int64 | float64] struct {
	processor numberprocessor.NumberProcessorInterface[T]
}

// Update sets the gauge value.
func (g *UpdateGaugeValueStrategyEngine[T]) Update(_, newValue string) (string, error) {
	// Парсим новое значение
	new, err := g.processor.Parse(newValue)
	if err != nil {
		return "", ErrUnprocessableValue
	}
	return g.processor.Format(new), nil
}
