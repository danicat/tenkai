package main

import (
	"log"
	"net/http"

	"experiment/internal/api"
	"experiment/internal/service"
)

func main() {
	svc := service.NewService()
	api := api.NewAPI(svc)

	http.HandleFunc("/users", api.HandleUsers)
	http.HandleFunc("/users/", api.HandleUserByID)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
