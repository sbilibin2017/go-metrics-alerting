package updatevalue

// StrategyUpdateEngineInterface defines the behavior for metric update strategies.
type UpdateValueStrategyEngineInterface interface {
	Update(currentValue, newValue string) (string, error)
}
