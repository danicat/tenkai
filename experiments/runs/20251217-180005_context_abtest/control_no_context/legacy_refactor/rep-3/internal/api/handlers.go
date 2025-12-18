package api

import (
	"encoding/json"
	"net/http"

	"experiment/internal/service"
)

// UserHandler provides HTTP handlers for user-related requests
type UserHandler struct {
	Service *service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

// GetUsersHandler handles GET requests for all users and POST requests to create a user
func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetUsers(w, r)
	case http.MethodPost:
		h.CreateUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetUserByIDHandler handles GET requests for a single user by ID
func (h *UserHandler) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]

	switch r.Method {
	case http.MethodGet:
		h.GetUserByID(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetUsers retrieves all users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := h.Service.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser := h.Service.CreateUser(u)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// GetUserByID retrieves a user by ID
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request, id string) {
	u, ok := h.Service.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}
