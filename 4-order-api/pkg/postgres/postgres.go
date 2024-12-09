package postgres

import (
	"advancedGo/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type db struct {
	*gorm.DB
}

func New(c *configs.Config) *db {
	database, err := gorm.Open(postgres.Open(c.DB.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &db{database}
}
