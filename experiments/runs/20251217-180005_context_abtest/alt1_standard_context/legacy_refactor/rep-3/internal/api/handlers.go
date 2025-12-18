package api

import (
	"encoding/json"
	"net/http"

	"rep-3/internal/service"
)

// GetUsers handles the GET request for all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	list := service.GetAllUsers()
	json.NewEncoder(w).Encode(list)
}

// CreateUser handles the POST request to create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u := service.CreateUser(req.Name, req.Email)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

// GetUser handles the GET request for a single user by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	u, ok := service.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}
