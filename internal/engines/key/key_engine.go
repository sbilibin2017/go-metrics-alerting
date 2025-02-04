package key

import (
	"strings"
)

// KeyEngine реализует движок
type KeyEngine struct{}

// NewKeyEngine создает новый движок.
func NewKeyEngine() *KeyEngine {
	return &KeyEngine{}
}

// Encode комбинирует тип и название метрики "type:name".
func (k *KeyEngine) Encode(metricType, metricName string) string {
	if metricType == "" || metricName == "" {
		return ""
	}
	return metricType + ":" + metricName
}

// Decode разбивает закодированный ключ на тип и название метрики.
func (k *KeyEngine) Decode(key string) (string, string, error) {
	metricType, metricName, found := strings.Cut(key, ":")
	if !found || metricType == "" || metricName == "" {
		return "", "", ErrInvalidKeyFormat
	}
	return metricType, metricName, nil
}

// Проверка соответствия интерфейсу
var _ KeyEngineInterface = (*KeyEngine)(nil)
