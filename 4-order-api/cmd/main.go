package main

import (
	"advancedGo/configs"
	"advancedGo/internal/auth"
	"advancedGo/internal/order"
	"advancedGo/internal/product"
	"advancedGo/internal/user"
	"advancedGo/pkg/db"
	"advancedGo/pkg/middleware"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func App(config *configs.Config) http.Handler {
	db := db.New(config)

	router := http.NewServeMux()

	// Repository
	productRepository := product.NewProductRepository(db)
	userRepository := user.NewUserRepository(db)
	orderRepository := order.NewOrderRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)

	// Handler
	auth.NewHandler(router, auth.AuthHandlerDeps{
		Config:      config,
		AuthService: authService,
	})

	product.NewHandler(router, product.ProductHandlerDeps{
		ProductRepository: productRepository,
	})
	order.NewHandler(router, order.OrderHandlerDeps{
		Config:             config,
		OrderRepository:    orderRepository,
		IProductRepository: productRepository,
		IUserRepository:    userRepository,
	})

	// Middlewares
	var middlewaresChain = middleware.Chain(
		middleware.Logging,
	)

	return middlewaresChain(router)

}

func main() {

	// Logrus
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     true,
	})

	config := configs.LoadConfig()

	server := http.Server{
		Addr:    ":8081",
		Handler: App(config),
	}

	fmt.Println("Start listening on port 8081")
	server.ListenAndServe()
}
