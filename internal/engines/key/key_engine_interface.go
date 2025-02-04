package key

// KeyEngineInterface определяет методы для  работы с движком ключей.
type KeyEngineInterface interface {
	Encode(metricType, metricName string) string
	Decode(key string) (string, string, error)
}
