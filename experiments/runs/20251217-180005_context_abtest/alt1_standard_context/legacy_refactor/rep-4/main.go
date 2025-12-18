package main

import (
	"log"
	"net/http"

	"experiment/internal/api"
	"experiment/internal/service"
)

func main() {
	// Initialize service
	userService := service.NewUserService()

	// Initialize handlers with the service
	userHandlers := api.NewUserHandlers(userService)

	// Register HTTP routes
	http.HandleFunc("/users", userHandlers.HandleUsers)
	http.HandleFunc("/users/", userHandlers.HandleUserByID)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}