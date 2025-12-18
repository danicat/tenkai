package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"legacy_refactor/internal/service"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler creates and returns a new UserHandler
func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// HandleUsers manages requests to /users
func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllUsers(w, r)
	case http.MethodPost:
		h.CreateUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleUserByID manages requests to /users/{id}
func (h *UserHandler) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	switch r.Method {
	case http.MethodGet:
		h.GetUserByID(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetAllUsers handles GET requests for all users
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.service.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUser handles POST requests to create a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user service.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.service.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// GetUserByID handles GET requests for a user by ID
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request, id string) {
	user, ok := h.service.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
}
