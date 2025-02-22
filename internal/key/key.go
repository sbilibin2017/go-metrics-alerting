package key

import (
	"strings"
)

const sep string = ":"

// Encode принимает произвольное количество строк и кодирует их в одну строку с использованием Join
func Encode(values ...string) string {
	return strings.Join(values, sep)
}

// Decode принимает закодированную строку и декодирует её, возвращая срез строк
func Decode(key string) []string {
	return strings.Split(key, sep)
}
