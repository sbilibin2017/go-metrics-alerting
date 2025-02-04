package updatevalue

import (
	"fmt"
	"go-metrics-alerting/internal/engines/numberprocessor"
)

// Обработчик обновлений для измерителей (работает с float64)
type UpdateGaugeValueStrategyEngine struct {
	processor numberprocessor.NumberProcessorEngineInterface[float64]
}

func NewUpdateGaugeValueStrategyEngine(
	processor numberprocessor.NumberProcessorEngineInterface[float64],
) *UpdateGaugeValueStrategyEngine {
	return &UpdateGaugeValueStrategyEngine{processor: processor}
}

// Реализация метода Update для обновления значения измерителя
func (g *UpdateGaugeValueStrategyEngine) Update(_, newValue string) (string, error) {
	// Парсим новое значение
	new, err := g.processor.Parse(newValue)
	if err != nil {
		return "", fmt.Errorf("error parsing new value '%s': %w", newValue, err)
	}

	// Возвращаем отформатированное новое значение
	return g.processor.Format(new), nil
}
