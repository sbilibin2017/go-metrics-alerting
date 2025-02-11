package engines

import (
	"errors"
	"strings"
)

const (
	KeySeparator   string = ":"
	KeyEmptyString string = ""
)

var (
	ErrInvalidKeyFormat error = errors.New("invalid key format")
)

// KeyEngine — структура для работы с ключами
type KeyEngine struct{}

// Encode генерирует ключ для метрики, соединяя тип метрики и её имя через ":"
func (ke *KeyEngine) Encode(mt string, mn string) string {
	return mt + KeySeparator + mn
}

// Decode расшифровывает ключ, разделяя его на тип и имя метрики
func (ke *KeyEngine) Decode(key string) (string, string, error) {
	parts := strings.Split(key, KeySeparator)
	if len(parts) != 2 || parts[0] == KeyEmptyString || parts[1] == KeyEmptyString {
		return KeyEmptyString, KeyEmptyString, ErrInvalidKeyFormat
	}
	return parts[0], parts[1], nil
}
