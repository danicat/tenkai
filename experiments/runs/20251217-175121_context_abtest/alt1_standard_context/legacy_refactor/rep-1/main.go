package main

import (
	"log"
	"net/http"

	"legacy_refactor/internal/api"
	"legacy_refactor/internal/service"
)

func main() {
	// Initialize service
	userService := service.NewInMemoryUserService()

	// Initialize handlers with the service
	handlers := api.NewUserHandlers(userService)

	// Register routes
	handlers.RegisterRoutes()

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}