package main

import (
	"log"
	"net/http"

	"rep-1/internal/api"
	"rep-1/internal/service"
)

func main() {
	// Initialize service and handler layers
	userService := service.NewUserService()
	userHandler := api.NewUserHandler(userService)

	// Register HTTP routes
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUsersHandler(w, r)
		case http.MethodPost:
			userHandler.CreateUserHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUserByIDHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}