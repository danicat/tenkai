package main

import (
	"log"
	"net/http"

	"experiment/internal/api"
	"experiment/internal/service"
)

func main() {
	// Initialize the UserService
	userService := service.NewUserService()

	// Initialize the UserHandler with the UserService
	userHandler := api.NewUserHandler(userService)

	// Register the handlers
	http.HandleFunc("/users", userHandler.ServeHTTP)
	http.HandleFunc("/users/", userHandler.ServeHTTP) // This will be handled by the UserHandler itself for specific IDs

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}