package key

import (
	"fmt"
	"strings"
)

// KeyEngine реализует KeyEngineInterface.
type KeyEngine struct{}

// NewKeyEngine создает новый экземпляр KeyEngine.
func NewKeyEngine() KeyEngineInterface {
	return &KeyEngine{}
}

// Encode комбинирует тип и название метрики в строку "type:name".
func (k *KeyEngine) Encode(key *Key) (string, error) {
	if key == nil || key.MetricType == "" || key.MetricName == "" {
		return "", ErrInvalidKeyFormat
	}
	return fmt.Sprintf("%s:%s", key.MetricType, key.MetricName), nil
}

// Decode разбивает строку "type:name" на тип и название метрики.
func (k *KeyEngine) Decode(key string) (*Key, error) {
	metricType, metricName, found := strings.Cut(key, ":")
	if !found || metricType == "" || metricName == "" {
		return nil, ErrInvalidKeyFormat
	}
	return &Key{MetricType: metricType, MetricName: metricName}, nil
}
