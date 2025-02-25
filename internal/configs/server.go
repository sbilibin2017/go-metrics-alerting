package configs

import "time"

// Структура конфигурации для сервера
type ServerConfig struct {
	Address         string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval   time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH" envDefault:"data.dump"`
	Restore         bool          `env:"RESTORE" envDefault:"false"`
}
