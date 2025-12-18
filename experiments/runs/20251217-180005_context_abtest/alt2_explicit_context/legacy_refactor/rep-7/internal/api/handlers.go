package api

import (
	"encoding/json"
	"net/http"

	"rep-7/internal/service"
)

// GetUsersHandler handles the GET /users request
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := service.GetUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUserHandler handles the POST /users request
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser := service.CreateUser(u)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// GetUserHandler handles the GET /users/{id} request
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	u, ok := service.GetUser(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}
