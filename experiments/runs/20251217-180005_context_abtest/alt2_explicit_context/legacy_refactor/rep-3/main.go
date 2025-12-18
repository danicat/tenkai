package main

import (
	"log"
	"net/http"
	"rep-3/internal/api"
	"rep-3/internal/service"
)

func main() {
	userService := service.NewInMemoryUserService()
	userHandler := api.NewUserHandler(userService)

	userHandler.SetupRoutes()

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
