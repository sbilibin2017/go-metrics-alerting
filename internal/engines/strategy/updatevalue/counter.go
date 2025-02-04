package updatevalue

import (
	"go-metrics-alerting/internal/engines/numberprocessor"
)

// CounterUpdateStrategyEngine handles counter metrics.
type UpdateCounterValueStrategyEngine[T int64 | float64] struct {
	processor numberprocessor.NumberProcessorInterface[T]
}

// Update increments the current counter value.
func (c *UpdateCounterValueStrategyEngine[T]) Update(currentValue, newValue string) (string, error) {
	current, err := c.processor.Parse(currentValue)
	if err != nil {
		return "", ErrUnprocessableValue
	}

	new, err := c.processor.Parse(newValue)
	if err != nil {
		return "", ErrUnprocessableValue
	}

	// Сложение значений и возвращение отформатированного результата.
	return c.processor.Format(current + new), nil
}
