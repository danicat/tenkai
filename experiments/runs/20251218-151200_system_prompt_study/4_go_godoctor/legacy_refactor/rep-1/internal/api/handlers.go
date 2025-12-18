package api

import (
	"encoding/json"
	"net/http"

	"experiment/internal/service"
)

// Handlers struct holds the service dependency
type Handlers struct {
	UserService *service.UserService
}

// GetUsersHandler handles GET requests for /users
func (h *Handlers) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := h.UserService.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUserHandler handles POST requests for /users
func (h *Handlers) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user service.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.UserService.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// GetUserByIDHandler handles GET requests for /users/{id}
func (h *Handlers) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	user, ok := h.UserService.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
}
