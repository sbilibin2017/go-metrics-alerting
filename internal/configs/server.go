package configs

import "time"

// ServerConfig now embeds the above configurations
type ServerConfig struct {
	Address         string
	DatabaseDSN     string
	FileStoragePath string
	StoreInterval   time.Duration
	Restore         bool
}

func NewServerConfig(
	address string,
	databaseDSN string,
	fileStoragePath string,
	storeInterval time.Duration,
	restore bool,
) *ServerConfig {
	return &ServerConfig{
		Address:         address,
		DatabaseDSN:     databaseDSN,
		FileStoragePath: fileStoragePath,
		StoreInterval:   storeInterval,
		Restore:         restore,
	}
}
