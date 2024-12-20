package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   dbConfig
	Auth authConfig
}

type dbConfig struct {
	DSN string
}

type authConfig struct {
	Token string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		DB: dbConfig{
			DSN: os.Getenv("DSN"),
		},
		Auth: authConfig{
			Token: os.Getenv("TOKEN"),
		},
	}

}
