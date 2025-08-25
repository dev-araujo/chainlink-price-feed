package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RpcURL     string
	ServerPort string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Não foi possível carregar o arquivo .env")
	}

	return &Config{
		RpcURL:     os.Getenv("SEPOLIA_RPC_URL"),
		ServerPort: os.Getenv("SERVER_PORT"),
	}
}
