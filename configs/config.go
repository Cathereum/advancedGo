package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db   dbConfig
	Auth authConfig
}

type dbConfig struct {
	Dsn string
}

type authConfig struct {
	Token string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}

	return &Config{
		Db: dbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: authConfig{
			Token: os.Getenv("TOKEN"),
		},
	}
}
