package main

import (
	"log"
	"net/http"

	"rep-1/internal/api"
	"rep-1/internal/service"
)

func main() {
	// Initialize service and handlers
	userService := service.NewUserService()
	userHandler := api.NewUserHandler(userService)

	// Register HTTP routes
	http.HandleFunc("/users", userHandler.HandleUsers)
	http.HandleFunc("/users/", userHandler.HandleUserByID)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}