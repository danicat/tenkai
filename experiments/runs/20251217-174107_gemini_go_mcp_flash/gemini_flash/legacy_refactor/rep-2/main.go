package main

import (
	"log"
	"net/http"

	"legacy_refactor/rep-2/internal/api"
	"legacy_refactor/rep-2/internal/service"
)

func main() {
	// Initialize service
	userService := service.NewUserService()

	// Initialize handlers with the service
	handlers := api.NewUserHandlers(userService)

	// Register HTTP routes
	http.HandleFunc("/users", handlers.HandleUsers)
	http.HandleFunc("/users/", handlers.HandleUserByID)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}