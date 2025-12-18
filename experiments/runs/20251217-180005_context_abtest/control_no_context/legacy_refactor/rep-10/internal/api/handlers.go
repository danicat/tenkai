package api

import (
	"encoding/json"
	"net/http"

	"experiment/internal/service"
)

// GetUsersHandler handles GET requests for users
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	list := service.GetAllUsers()
	json.NewEncoder(w).Encode(list)
}

// CreateUserHandler handles POST requests for users
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := service.CreateUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// GetUserByIDHandler handles GET requests for a user by ID
func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	u, ok := service.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}
