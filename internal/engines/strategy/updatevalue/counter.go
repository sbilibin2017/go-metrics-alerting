package updatevalue

import (
	"fmt"
	"go-metrics-alerting/internal/engines/numberprocessor"
)

// Обработчик обновлений для счетчиков
type UpdateCounterValueStrategyEngine struct {
	processor numberprocessor.NumberProcessorEngineInterface[int64]
}

func NewUpdateCounterValueStrategyEngine(
	processor numberprocessor.NumberProcessorEngineInterface[int64],
) *UpdateCounterValueStrategyEngine {
	return &UpdateCounterValueStrategyEngine{processor: processor}
}

// Реализация метода Update для обновления значения счетчика
func (c *UpdateCounterValueStrategyEngine) Update(currentValue, newValue string) (string, error) {
	// Парсим текущее значение
	current, err := c.processor.Parse(currentValue)
	if err != nil {
		return "", fmt.Errorf("error parsing current value '%s': %w", currentValue, err)
	}

	// Парсим новое значение
	new, err := c.processor.Parse(newValue)
	if err != nil {
		return "", fmt.Errorf("error parsing new value '%s': %w", newValue, err)
	}

	// Возвращаем отформатированное обновленное значение
	return c.processor.Format(current + new), nil
}
