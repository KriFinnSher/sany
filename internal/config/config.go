package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost     string
	ServerPort     string
	DataSourcePath string
}

func init() {
	godotenv.Load()
}

// MustLoad reads required environment configuration or panics when it is incomplete.
func MustLoad() *Config {
	sHost := os.Getenv("SERVER_HOST")
	sPort := os.Getenv("SERVER_PORT")
	dsPath := os.Getenv("DATASOURCE_PATH")

	if sHost == "" || sPort == "" || dsPath == "" {
		panic("failed to load config")
	}

	return &Config{
		ServerHost:     sHost,
		ServerPort:     sPort,
		DataSourcePath: dsPath,
	}
}
