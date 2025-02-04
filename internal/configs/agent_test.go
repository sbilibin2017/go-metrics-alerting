package configs

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Тестируем LoadAgentConfigFromEnv
func TestLoadAgentConfigFromEnv(t *testing.T) {
	// Создаем временный .env файл
	file, err := os.Create(".env")
	if err != nil {
		t.Fatalf("Error creating .env file: %v", err)
	}
	defer file.Close()

	// Записываем переменные окружения в .env файл
	_, err = file.WriteString(`
		ADDRESS=http://testaddress:8080
		REPORT_INTERVAL=5s
		POLL_INTERVAL=1s
	`)
	if err != nil {
		t.Fatalf("Error writing to .env file: %v", err)
	}

	// Загружаем конфигурацию из .env
	err = godotenv.Load() // Этот вызов загрузит файл .env
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	config, err := LoadAgentConfigFromEnv()
	if err != nil {
		t.Fatalf("Error loading config from .env: %v", err)
	}

	// Проверяем, что конфигурация загружена правильно
	assert.Equal(t, "http://testaddress:8080", config.Address)
	assert.Equal(t, 5*time.Second, config.ReportInterval)
	assert.Equal(t, 1*time.Second, config.PollInterval)

	// Удаляем временный .env файл
	os.Remove(".env")
}

// Тестируем LoadAgentConfigFromEnv с дефолтными значениями, если .env отсутствует
func TestLoadAgentConfigFromEnv_Defaults(t *testing.T) {
	// Удаляем .env файл, если он существует
	os.Remove(".env")

	// Загружаем конфигурацию из окружения (без .env файла)
	config, err := LoadAgentConfigFromEnv()

	// Проверяем, что функция вернула nil (так как .env не существует)
	assert.Nil(t, config)
	assert.Nil(t, err)
}

// Тестируем LoadAgentConfigFromFlags
func TestLoadAgentConfigFromFlags(t *testing.T) {
	// Создаем флаги для конфигурации
	os.Args = []string{"", "-a", "http://testaddress:8081", "-r", "15s", "-p", "3s"}

	// Загружаем конфигурацию из флагов
	config, err := LoadAgentConfigFromFlags()
	if err != nil {
		t.Fatalf("Error loading config from flags: %v", err)
	}

	// Проверяем, что флаги правильно обработаны
	assert.Equal(t, "http://testaddress:8081", config.Address)
	assert.Equal(t, 15*time.Second, config.ReportInterval)
	assert.Equal(t, 3*time.Second, config.PollInterval)
}

// Тестируем LoadAgentConfigFromFlags с дефолтными значениями
func TestLoadAgentConfigFromFlags_Defaults(t *testing.T) {
	// Устанавливаем флаги по умолчанию
	os.Args = []string{"", "-a", "localhost:8080"}

	// Загружаем конфигурацию из флагов
	config, err := LoadAgentConfigFromFlags()
	if err != nil {
		t.Fatalf("Error loading config from flags: %v", err)
	}

	// Проверяем, что дефолтные значения применены
	assert.Equal(t, "localhost:8080", config.Address)
	assert.Equal(t, 10*time.Second, config.ReportInterval)
	assert.Equal(t, 2*time.Second, config.PollInterval)
}
