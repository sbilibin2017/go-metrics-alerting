package updatevalue

// GaugeUpdateStrategyEngine handles gauge metrics.
type UpdateGaugeValueStrategyEngine struct{}

// Update sets the gauge value.
func (g *UpdateGaugeValueStrategyEngine) Update(_, newValue string) (string, error) {
	new, err := parseNumber[float64](newValue)
	if err != nil {
		return "", err
	}

	return formatNumber[float64](new), nil
}
