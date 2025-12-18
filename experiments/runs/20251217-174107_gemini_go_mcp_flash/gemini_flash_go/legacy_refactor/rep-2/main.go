package main

import (
	"log"
	"net/http"
	"legacy_refactor/rep-2/internal/api"
	"legacy_refactor/rep-2/internal/service"
)

func main() {
	// Initialize service layer
	userService := service.NewUserService()

	// Initialize API handlers with the service
	apiHandlers := api.NewHandlers(userService)

	// Register HTTP routes
	http.HandleFunc("/users", apiHandlers.HandleUsers)
	http.HandleFunc("/users/", apiHandlers.HandleUserByID)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}