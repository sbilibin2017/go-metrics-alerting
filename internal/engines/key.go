package engines

import (
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
	"strings"
)

// KeyEngine — структура для работы с ключами
type KeyEngine struct{}

// Encode генерирует ключ для метрики, соединяя тип метрики и её имя через ":"
func (ke *KeyEngine) Encode(mt string, mn string) string {
	// Используем константу для разделителя
	return mt + types.KeySeparator + mn
}

// Decode расшифровывает ключ, разделяя его на тип и имя метрики
func (ke *KeyEngine) Decode(key string) (string, string, error) {
	// Используем константу для разделителя
	parts := strings.Split(key, types.KeySeparator)
	if len(parts) != 2 || parts[0] == types.EmptyString || parts[1] == types.EmptyString {
		return types.EmptyString, types.EmptyString, errors.ErrInvalidKeyFormat
	}

	// Возвращаем строковые значения без использования доменных структур
	return parts[0], parts[1], nil
}
