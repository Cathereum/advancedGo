package main

import (
	"advancedGo/configs"
	"advancedGo/internal/auth"
	"fmt"
	"net/http"
)

func main() {
	config := configs.LoadConfig()

	router := http.NewServeMux()
	auth.NewHandler(router, auth.AuthHandlerDeps{
		Config: config,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Start listening on port 8081")
	server.ListenAndServe()
}
