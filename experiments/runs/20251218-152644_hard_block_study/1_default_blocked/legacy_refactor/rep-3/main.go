package main

import (
	"log"
	"net/http"

	"experiment/internal/api"
	"experiment/internal/service"
)

func main() {
	// Initialize the user service
	userService := service.NewUserService()

	// Initialize user handlers with the service
	userHandlers := api.NewUserHandlers(userService)

	// Register HTTP routes
	userHandlers.RegisterRoutes()

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
