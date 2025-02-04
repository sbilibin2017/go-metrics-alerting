package updatevalue

// CounterUpdateStrategyEngine handles counter metrics.
type UpdateCounterValueStrategyEngine struct{}

// Update increments the current counter value.
func (c *UpdateCounterValueStrategyEngine) Update(currentValue, newValue string) (string, error) {
	current, err := parseNumber[int64](currentValue)
	if err != nil {
		return "", err
	}

	new, err := parseNumber[int64](newValue)
	if err != nil {
		return "", err
	}

	return formatNumber[int64](current + new), nil
}
