package main

import (
	"log"
	"net/http"

	"legacy_refactor/internal/api"
	"legacy_refactor/internal/service"
)

func main() {
	svc := service.NewUserService()
	handler := api.NewUserHandler(svc)

	http.HandleFunc("/users", handler.HandleUsers)
	http.HandleFunc("/users/", handler.HandleUserByID)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
