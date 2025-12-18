package api

import (
	"encoding/json"
	"net/http"

	"rep-1/internal/service"
)

// ListUsersHandler handles listing all users
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := service.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUserHandler handles creating a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u = service.CreateUser(u)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

// GetUserByIDHandler handles retrieving a user by ID
func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	u, ok := service.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}
