package keyprocessor

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

// KeyDecoder структура для декодирования ключей
type KeyDecoder struct{}

func NewKeyDecoder() *KeyDecoder {
	return &KeyDecoder{}
}

// Метод Decode для декодирования строки в id и mtype с использованием strings.Split
func (d *KeyDecoder) Decode(key string) (id string, mtype string, ok bool) {
	parts := strings.Split(key, sep)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}
