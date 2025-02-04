package configs

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

type ServerConfig struct {
	Address string
}

func LoadServerConfigFromEnv() (*ServerConfig, error) {
	err := godotenv.Load()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	address := os.Getenv("ADDRESS")
	return &ServerConfig{
		Address: address,
	}, nil
}

func LoadServerConfigFromFlags() (*ServerConfig, error) {
	var config ServerConfig
	var rootCmd = &cobra.Command{
		Use:   "server",
		Short: "Server for handling HTTP requests",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	rootCmd.Flags().StringVarP(&config.Address, "address", "a", "localhost:8080", "HTTP server endpoint")
	if err := rootCmd.Execute(); err != nil {
		return nil, err
	}
	return &config, nil
}
