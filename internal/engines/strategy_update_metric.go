package engines

import (
	"strconv"
)

// Интерфейс для стратегий работы с метриками
type StrategyUpdateEngineInterface interface {
	Update(currentValue string, newValue string) (string, error)
}

// CounterStrategy is the strategy for handling counter metrics
type CounterUpdateStrategyEngine struct{}

func NewCounterUpdateStrategyEngine() *CounterUpdateStrategyEngine {
	return &CounterUpdateStrategyEngine{}
}

// Update increments the current counter value by the new value
// Теперь метод принимает строковые значения и возвращает строковое значение.
func (c *CounterUpdateStrategyEngine) Update(currentValue string, newValue string) (string, error) {
	// Parse the increment value (newValue)
	increment, err := ParseInt(newValue)
	if err != nil {
		return "", err // Invalid increment value
	}

	// Parse the current value
	current, err := ParseInt(currentValue)
	if err != nil {
		return "", err // Invalid current value
	}

	// Calculate the updated value without any restrictions
	updatedValue := current + increment

	// Return the updated value as a string
	return FormatInt(updatedValue), nil
}

func FormatInt(value int64) string {
	return strconv.FormatInt(value, 10)
}

func ParseInt(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

// GaugeStrategy is the strategy for handling gauge metrics
type GaugeUpdateStrategyEngine struct{}

func NewGaugeUpdateStrategyEngine() *GaugeUpdateStrategyEngine {
	return &GaugeUpdateStrategyEngine{}
}

// Update sets the current gauge value to the new value
// Теперь метод принимает строковые значения и возвращает строковое значение.
func (g *GaugeUpdateStrategyEngine) Update(currentValue string, newValue string) (string, error) {
	value, err := ParseFloat(newValue)
	if err != nil {
		return "", err
	}

	return FormatFloat(value), nil
}

func FormatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func ParseFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

// Ensure both strategies implement StrategyEngineInterface
var _ StrategyUpdateEngineInterface = &GaugeUpdateStrategyEngine{}
var _ StrategyUpdateEngineInterface = &CounterUpdateStrategyEngine{}
