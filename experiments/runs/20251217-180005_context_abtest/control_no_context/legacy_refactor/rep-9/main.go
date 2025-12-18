package main

import (
	"log"
	"net/http"

	"experiment/internal/api"
	"experiment/internal/service"
)

func main() {
	// Initialize service (if it had any stateful dependencies, they would be passed here)
	_ = service.ListUsers // Call a service function to ensure it's imported and initialized.

	http.HandleFunc("/users", api.HandleUsers)
	http.HandleFunc("/users/", api.HandleUserByID)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}