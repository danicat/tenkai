package api

import (
	"encoding/json"
	"net/http"

	"experiment/internal/service"
)

// ListUsersHandler handles GET requests to list all users
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := service.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUserHandler handles POST requests to create a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser := service.CreateUser(u)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// GetUserHandler handles GET requests for a specific user by ID
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	u, ok := service.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}
