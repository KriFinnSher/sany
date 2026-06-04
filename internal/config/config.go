package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost string
	ServerPort string
}

func init() {
	godotenv.Load()
}

func MustLoad() *Config {
	sHost := os.Getenv("SERVER_HOST")
	sPort := os.Getenv("SERVER_PORT")

	if sHost == "" || sPort == "" {
		panic("failed to load config")
	}

	return &Config{
		ServerHost: sHost,
		ServerPort: sPort,
	}
}
