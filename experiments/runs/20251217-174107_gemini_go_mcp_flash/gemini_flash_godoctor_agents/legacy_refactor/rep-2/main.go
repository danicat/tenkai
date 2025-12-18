package main

import (
	"log"
	"net/http"

	"rep-2/internal/api"
)

func main() {
	handlers := api.NewUserHandlers()

	http.HandleFunc("/users", handlers.HandleUsers)
	http.HandleFunc("/users/", handlers.HandleUserByID)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
