package main

import (
	"log"
	"net/http"

	"rep-2/internal/api"
	"rep-2/internal/service"
)

func main() {
	// Initialize service
	userService := service.NewUserService()

	// Initialize API handlers with the service
	handler := api.NewUserHandler(userService)

	// Register API routes
	handler.RegisterRoutes()

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}