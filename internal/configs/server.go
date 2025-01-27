package configs

import (
	"flag"
)

// Интерфейс для конфигурации сервера
type ServerConfigInterface interface {
	GetServerUrl() string
}

// Структура конфигурации для сервера
type ServerConfig struct {
	ServerUrl string // Адрес эндпоинта HTTP-сервера
}

// Функция для создания конфигурации сервера с флагами и переменными окружения
func NewServerConfig() *ServerConfig {
	// Флаг для адреса сервера
	serverURL := flag.String("a", "", "Адрес эндпоинта HTTP-сервера")
	flag.Parse()

	// Получаем значение из переменной окружения SERVER_URL или из флага
	finalServerURL := GetEnvOrDefault(*serverURL, "SERVER_URL", "localhost:8080")

	return &ServerConfig{
		ServerUrl: finalServerURL,
	}
}

// Реализация метода интерфейса ServerConfigInterface для ServerConfig
func (config *ServerConfig) GetServerUrl() string {
	return config.ServerUrl
}
