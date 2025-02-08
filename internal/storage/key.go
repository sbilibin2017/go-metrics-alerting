package storage

import (
	"errors"
	"strings"
)

const (
	sep                   = ":"
	emptyKeyString string = ""
)

// ErrInvalidKeyFormat — ошибка, если ключ не имеет правильного формата
var ErrInvalidKeyFormat = errors.New("invalid key format")

// KeyManager — структура для работы с ключами
type keyProcessor struct{}

// NewKeyManager — создание нового KeyManager
func NewKeyProcessor() *keyProcessor {
	return &keyProcessor{}
}

// Encode генерирует ключ для метрики, соединяя тип метрики и её имя через ":"
func (km *keyProcessor) Encode(metricType string, metricName string) string {
	return strings.Join([]string{metricType, metricName}, sep)
}

// Decode расшифровывает ключ, разделяя его на тип и имя метрики
func (km *keyProcessor) Decode(key string) (string, string, error) {
	parts := strings.Split(key, sep)
	if len(parts) != 2 || parts[0] == emptyKeyString || parts[1] == emptyKeyString {
		return emptyKeyString, emptyKeyString, ErrInvalidKeyFormat
	}
	return parts[0], parts[1], nil
}
