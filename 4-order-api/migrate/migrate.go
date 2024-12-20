package main

import (
	"advancedGo/internal/product"
	"advancedGo/internal/user"
	"os"

	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		panic(err)
	}

	database, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&product.Product{}, &user.User{})
}
