package configs

import (
	"github.com/caarlos0/env/v6"
)

// Структура конфигурации для сервера
type ServerConfig struct {
	Address string `env:"ADDRESS"` // Адрес сервера без default
}

// Функция для загрузки конфигурации сервера
func LoadServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	// Если переменная окружения не задана, устанавливаем значение по умолчанию
	if cfg.Address == "" {
		cfg.Address = ":8080"
	}

	return cfg, nil
}
