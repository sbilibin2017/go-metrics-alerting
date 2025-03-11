package helpers

import (
	"os"
	"time"
)

// GetStringFromEnv получает строковое значение из переменной окружения.
// Если переменная окружения не найдена, возвращает пустую строку.
func GetStringFromEnv(key string) *string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return nil
	}
	return &value
}

// GetBoolFromEnv получает булевое значение из переменной окружения.
// Если переменная окружения не найдена или её значение не "true", возвращает false.
func GetBoolFromEnv(key string) *bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return nil
	}
	v := value == "true"
	return &v
}

// GetDurationFromEnv получает продолжительность из переменной окружения.
// Если переменная окружения не найдена или значение не может быть преобразовано в time.Duration,
// возвращает дефолтное значение.
func GetDurationFromEnv(key string) *time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return nil
	}
	parsedValue, err := time.ParseDuration(value)
	if err != nil {
		return nil
	}
	return &parsedValue
}
