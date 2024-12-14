package main

import (
	"advancedGo/configs"
	"advancedGo/internal/product"
	"advancedGo/pkg/db"
	"advancedGo/pkg/middleware"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {

	// Logrus
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     true,
	})

	c := configs.LoadConfig()
	db := db.New(c)

	router := http.NewServeMux()

	// Repository
	productRepository := product.NewProductRepository(db)

	// Handler
	product.NewHandler(router, product.ProductHandlerDeps{
		ProductRepository: productRepository,
	})

	// Middlewares
	var middlewaresChain = middleware.Chain(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: middlewaresChain(router),
	}

	fmt.Println("Start listening on port 8081")
	server.ListenAndServe()
}
