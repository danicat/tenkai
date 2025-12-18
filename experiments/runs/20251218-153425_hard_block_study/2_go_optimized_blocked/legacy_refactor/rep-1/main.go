package main

import (
	"log"
	"net/http"

	"legacy_refactor/rep-1/internal/api"
	"legacy_refactor/rep-1/internal/service"
)

func main() {
	userService := service.NewUserService()
	userHandler := api.NewUserHandler(userService)

	http.HandleFunc("/users", userHandler.HandleUsers)
	http.HandleFunc("/users/", userHandler.HandleUserByID)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
