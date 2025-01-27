package engines

import (
	"errors"
	"fmt"
	"strings"
)

// KeyEngineInterface defines the methods for encoding and decoding keys
type KeyEngineInterface interface {
	// Encode combines the type and name into a single string "type:name"
	Encode(metricType string, metricName string) string

	// Decode splits the encoded key into its type and name parts
	Decode(key string) (string, string, error)
}

// KeyEngine implements the storage of keys
type KeyEngine struct{}

// NewKeyEngine creates a new instance of KeyEngine
func NewKeyEngine() *KeyEngine {
	return &KeyEngine{}
}

// Encode combines the type and name into a single string "type:name"
func (k *KeyEngine) Encode(metricType string, metricName string) string {
	return fmt.Sprintf("%s:%s", metricType, metricName)
}

// Decode splits the encoded key into its type and name parts
func (k *KeyEngine) Decode(key string) (string, string, error) {
	parts := strings.Split(key, ":")
	if len(parts) != 2 {
		return "", "", errors.New("invalid key format, expected 'type:name'")
	}
	return parts[0], parts[1], nil
}

var _ KeyEngineInterface = &KeyEngine{}
