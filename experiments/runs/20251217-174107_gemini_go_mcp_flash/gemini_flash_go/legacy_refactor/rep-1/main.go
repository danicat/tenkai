package main

import (
	"log"
	"net/http"

	"experiment/internal/api"
	"experiment/internal/service"
)

func main() {
	// Initialize the service layer
	svc := service.NewService()

	// Initialize the API handlers with the service
	apiHandlers := api.NewAPI(svc)

	// Register HTTP handlers
	http.HandleFunc("/users", apiHandlers.HandleUsers)
	http.HandleFunc("/users/", apiHandlers.HandleUserByID)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}