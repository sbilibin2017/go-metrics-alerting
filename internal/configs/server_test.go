package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тестируем LoadServerConfigFromEnv с дефолтным значением, если .env отсутствует
func TestLoadServerConfigFromEnv_Defaults(t *testing.T) {
	// Удаляем .env файл, если он существует
	os.Remove(".env")

	// Загружаем конфигурацию из окружения (без .env файла)
	config, err := LoadServerConfigFromEnv()

	// Проверяем, что функция вернула nil (так как .env не существует)
	assert.Nil(t, config)
	assert.Nil(t, err)
}

// Тестируем LoadServerConfigFromEnv с значениями из .env
func TestLoadServerConfigFromEnv_Values(t *testing.T) {
	// Создаем временный .env файл
	file, err := os.Create(".env")
	if err != nil {
		t.Fatalf("Error creating .env file: %v", err)
	}
	defer file.Close()

	// Записываем переменные окружения в .env файл
	_, err = file.WriteString(`
		ADDRESS=http://testaddress:8080
	`)
	if err != nil {
		t.Fatalf("Error writing to .env file: %v", err)
	}

	// Загружаем конфигурацию из .env
	config, err := LoadServerConfigFromEnv()
	if err != nil {
		t.Fatalf("Error loading config from .env: %v", err)
	}

	// Проверяем, что конфигурация загружена правильно
	assert.Equal(t, "http://testaddress:8080", config.Address)

	// Удаляем временный .env файл
	os.Remove(".env")
}

// Тестируем LoadServerConfigFromFlags с пользовательскими значениями флагов
func TestLoadServerConfigFromFlags_CustomFlags(t *testing.T) {
	// Устанавливаем флаги для теста
	os.Args = []string{"", "-a", "http://customaddress:9090"}

	// Загружаем конфигурацию из флагов
	config, err := LoadServerConfigFromFlags()
	if err != nil {
		t.Fatalf("Error loading config from flags: %v", err)
	}

	// Проверяем, что флаги правильно обработаны
	assert.Equal(t, "http://customaddress:9090", config.Address)
}

// Тестируем LoadServerConfigFromFlags с дефолтными значениями
func TestLoadServerConfigFromFlags_DefaultFlags(t *testing.T) {
	// Устанавливаем флаги по умолчанию
	os.Args = []string{"", "-a", "localhost:8080"}

	// Загружаем конфигурацию из флагов
	config, err := LoadServerConfigFromFlags()
	if err != nil {
		t.Fatalf("Error loading config from flags: %v", err)
	}

	// Проверяем, что дефолтное значение для адреса используется
	assert.Equal(t, "localhost:8080", config.Address)
}

// Тестируем LoadServerConfigFromFlags без указания флагов
func TestLoadServerConfigFromFlags_NoFlags(t *testing.T) {
	// Устанавливаем флаги по умолчанию
	os.Args = []string{"", "-a", "localhost:8080"}

	// Загружаем конфигурацию из флагов
	config, err := LoadServerConfigFromFlags()
	if err != nil {
		t.Fatalf("Error loading config from flags: %v", err)
	}

	// Проверяем, что дефолтное значение для адреса используется
	assert.Equal(t, "localhost:8080", config.Address)
}
