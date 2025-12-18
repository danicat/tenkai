package main

import (
	"log"
	"net/http"

	"experiment/internal/api"
)

func main() {
	api.RegisterRoutes()

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}