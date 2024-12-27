package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   DBConfig
	Auth AuthConfig
}

type DBConfig struct {
	DSN string
}

type AuthConfig struct {
	Token string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		DB: DBConfig{
			DSN: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			Token: os.Getenv("TOKEN"),
		},
	}

}
