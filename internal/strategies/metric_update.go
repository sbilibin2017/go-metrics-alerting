package strategies

import (
	"errors"
	"go-metrics-alerting/internal/types"
	"strconv"
)

// Ошибка, если значение Gauge не корректно
var (
	ErrInvalidGaugeValue    = errors.New("invalid gauge value")
	ErrInvalidCounterFormat = errors.New("invalid counter format")
	ErrSaveCounter          = errors.New("counter is not saved")
	ErrSaveGauge            = errors.New("gauge is not saved")
)

// Setter интерфейс для записи данных.
type Saver interface {
	Save(key string, value string) error
}

// GaugeUpdateStrategy - стратегия для обновления метрики типа gauge
type GaugeUpdateStrategy struct {
	saver Saver // Внедрение зависимостей через Setter
}

// Update для типа gauge перезаписывает старое значение новым
func (g *GaugeUpdateStrategy) Update(req *types.MetricsRequest, currentValue string) (*types.MetricsRequest, error) {
	// Проверяем, что значение Value присутствует
	if req.Value == nil {
		return nil, ErrInvalidGaugeValue
	}

	// Преобразуем значение в строку с помощью strconv
	updatedMetricValue := strconv.FormatFloat(*req.Value, 'f', -1, 64)

	// Сохраняем обновленную метрику с помощью Setter
	err := g.saver.Save(req.ID, updatedMetricValue)
	if err != nil {
		return nil, ErrSaveGauge
	}

	// Возвращаем обновленный запрос
	return req, nil
}

// CounterUpdateStrategy - стратегия для обновления метрики типа counter
type CounterUpdateStrategy struct {
	saver Saver // Внедрение зависимостей через Setter
}

// Update для типа counter добавляет новое значение к старому
func (c *CounterUpdateStrategy) Update(req *types.MetricsRequest, currentValue string) (*types.MetricsRequest, error) {
	intValue, err := strconv.ParseInt(currentValue, 10, 64)
	if err != nil {
		return nil, ErrInvalidCounterFormat
	}

	*req.Delta += intValue

	updatedMetricValue := strconv.FormatInt(*req.Delta, 10)

	err = c.saver.Save(req.ID, updatedMetricValue)
	if err != nil {
		return nil, ErrSaveCounter
	}

	// Возвращаем обновленный запрос
	return req, nil
}
