package keyencoder

import (
	"strings"
)

// Разделитель для ключа
const sep string = ":"

// KeyEncoder структура для кодирования ключей
type KeyEncoder struct{}

func NewKeyEncoder() *KeyEncoder {
	return &KeyEncoder{}
}

// Метод Encode для кодирования id и mtype в строку с использованием strings.Join
func (e *KeyEncoder) Encode(id string, mtype string) string {
	return strings.Join([]string{id, mtype}, sep)
}
