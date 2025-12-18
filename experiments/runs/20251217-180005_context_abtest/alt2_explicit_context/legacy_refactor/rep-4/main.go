package main

import (
	"log"
	"net/http"

	"rep-4/internal/api"
)

func main() {
	http.HandleFunc("/users", api.UsersHandler)
	http.HandleFunc("/users/", api.UserByIDHandler)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}