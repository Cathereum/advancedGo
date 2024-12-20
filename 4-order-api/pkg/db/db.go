package db

import (
	"advancedGo/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func New(c *configs.Config) *DB {
	db, err := gorm.Open(postgres.Open(c.DB.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &DB{db}
}
