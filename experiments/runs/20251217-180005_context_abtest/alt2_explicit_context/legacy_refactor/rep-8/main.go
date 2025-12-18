package main

import (
	"log"
	"net/http"

	"experiment/internal/api"
	"experiment/internal/service"
)

func main() {
	// Initialize the service
	userService := service.NewUserService()

	// Initialize the API handlers with the service
	userHandlers := api.NewUserHandlers(userService)

	// Register HTTP handlers
	http.HandleFunc("/users", userHandlers.HandleUsers)
	http.HandleFunc("/users/", userHandlers.HandleUserByID)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}