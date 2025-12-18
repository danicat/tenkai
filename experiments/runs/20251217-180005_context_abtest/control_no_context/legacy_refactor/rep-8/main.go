package main

import (
	"log"
	"net/http"

	"rep-8/internal/api"
	"rep-8/internal/service"
)

func main() {
	userService := service.NewUserService()
	handler := api.NewUserHandler(userService)

	http.HandleFunc("/users", handler.HandleUsers)
	http.HandleFunc("/users/", handler.HandleUserByID)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}