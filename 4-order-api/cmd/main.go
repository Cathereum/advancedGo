package main

import (
	"advancedGo/configs"
	"advancedGo/pkg/postgres"
)

func main() {
	c := configs.LoadConfig()
	postgres.New(c)
}
