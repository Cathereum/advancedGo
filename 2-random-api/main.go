package main

import (
	"fmt"
	"net/http"
)

func main() {

	router := http.NewServeMux()
	NewNumRandomHAndler(router)
	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Start listening on port 8081")
	server.ListenAndServe()
}
