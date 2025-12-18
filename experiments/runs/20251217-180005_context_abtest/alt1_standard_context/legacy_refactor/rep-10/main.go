package main

import (
	"log"
	"net/http"

	"legacy_refactor/rep-10/internal/api"
	"legacy_refactor/rep-10/internal/service"
)

func main() {
	// Initialize service layer
	userService := service.NewUserService()

	// Initialize API handlers
	userHandlers := api.NewUserHandlers(userService)

	// Register HTTP routes
	http.HandleFunc("/users", userHandlers.HandleUsers)
	http.HandleFunc("/users/", userHandlers.HandleUserByID)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}