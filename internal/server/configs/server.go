package configs

// Структура конфигурации для сервера
type ServerConfig struct {
	Address string `env:"ADDRESS"` // Адрес сервера без default
}
