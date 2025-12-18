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

	// Initialize API handlers with the user service
	handlers := api.NewUserHandlers(userService)

	// Register HTTP handlers
	http.HandleFunc("/users", handlers.UsersHandler)
	http.HandleFunc("/users/", handlers.UserByIDHandler)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
