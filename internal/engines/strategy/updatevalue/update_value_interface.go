package updatevalue

// Интерфейс для обновления значений стратегии.
type UpdateValueStrategyEngineInterface interface {
	// Update обновляет значение, используя текущие и новые значения.
	Update(currentValue, newValue string) (string, error)
}
