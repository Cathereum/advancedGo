package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB dbConfig
}

type dbConfig struct {
	DSN string
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
	}

}
