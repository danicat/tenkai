package api

import (
	"encoding/json"
	"net/http"

	"legacy_refactor/rep-2/internal/service"
)

// Handler handles HTTP requests
type Handler struct {
	service *service.Service
}

// NewHandler creates and returns a new Handler instance
func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

// ListUsersHandler handles GET /users requests
func (h *Handler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := h.service.ListUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUserHandler handles POST /users requests
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser := h.service.CreateUser(u)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// GetUserHandler handles GET /users/{id} requests
func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	u, ok := h.service.GetUser(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}
