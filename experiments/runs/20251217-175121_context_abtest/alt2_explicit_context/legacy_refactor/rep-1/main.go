package main

import (
	"log"
	"net/http"

	"experiment/internal/api"
	"experiment/internal/service"
)

func main() {
	// Initialize service
	svc := service.NewUserService()

	// Initialize handlers with the service
	handlers := api.NewUserHandlers(svc)

	// Register routes
	handlers.RegisterRoutes()

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}