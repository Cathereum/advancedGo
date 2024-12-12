package main

import (
	"advancedGo/configs"
	"advancedGo/internal/product"
	"advancedGo/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	c := configs.LoadConfig()
	db := db.New(c)

	router := http.NewServeMux()

	// Repository
	productRepository := product.NewProductRepository(db)

	// Handler
	product.NewHandler(router, product.ProductHandlerDeps{
		ProductRepository: productRepository,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Start listening on port 8081")
	server.ListenAndServe()
}
