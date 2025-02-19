package configs

// Структура конфигурации для сервера
type ServerConfig struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}
