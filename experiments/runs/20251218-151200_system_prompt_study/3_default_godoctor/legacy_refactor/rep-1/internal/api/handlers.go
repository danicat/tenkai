package api

import (
	"encoding/json"
	"net/http"

	"rep-1/internal/service"
)

// UserHandlers holds the dependencies for user HTTP handlers
type UserHandlers struct {
	Service service.UserService
}

// GetUsers handles GET requests for /users
func (h *UserHandlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := h.Service.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUser handles POST requests for /users
func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := h.Service.AddUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// GetUserByID handles GET requests for /users/{id}
func (h *UserHandlers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	u, ok := h.Service.FindUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}
