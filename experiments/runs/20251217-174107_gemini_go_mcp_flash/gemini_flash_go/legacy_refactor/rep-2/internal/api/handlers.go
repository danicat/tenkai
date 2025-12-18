package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"legacy_refactor/rep-2/internal/service"
)

// Handlers holds the dependencies for API handlers
type Handlers struct {
	UserService *service.UserService
}

// NewHandlers creates new API handlers
func NewHandlers(us *service.UserService) *Handlers {
	return &Handlers{
		UserService: us,
	}
}

// HandleUsers handles requests for /users
func (h *Handlers) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetUsers(w, r)
	case http.MethodPost:
		h.CreateUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetUsers returns all users
func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := h.UserService.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUser creates a new user
func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.UserService.CreateUser(u.Name, u.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// HandleUserByID handles requests for /users/{id}
func (h *Handlers) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	switch r.Method {
	case http.MethodGet:
		h.GetUserByID(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetUserByID returns a user by ID
func (h *Handlers) GetUserByID(w http.ResponseWriter, r *http.Request, id string) {
	u, err := h.UserService.GetUserByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}
