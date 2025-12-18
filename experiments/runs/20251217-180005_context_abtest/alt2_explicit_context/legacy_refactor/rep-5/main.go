package main

import (
	"log"
	"net/http"

	"rep-5/internal/api"
)

func main() {
	http.HandleFunc("/users", api.GetUsersHandler)
	http.HandleFunc("/users/", api.GetUserByIDHandler)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}