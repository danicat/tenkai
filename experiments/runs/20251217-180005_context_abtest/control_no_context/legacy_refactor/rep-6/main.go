package main

import (
	"log"
	"net/http"

	"experiment/internal/api"
	"experiment/internal/service"
)

func main() {
	userService := service.NewService()
	userHandlers := api.NewUserHandlers(userService)

	http.HandleFunc("/users", userHandlers.HandleUsers)
	http.HandleFunc("/users/", userHandlers.HandleUserByID)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}