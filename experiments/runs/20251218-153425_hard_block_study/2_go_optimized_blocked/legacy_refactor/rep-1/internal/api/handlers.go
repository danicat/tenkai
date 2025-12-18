package api

import (
	"encoding/json"
	"net/http"

	"legacy_refactor/rep-1/internal/service"
)

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	UserService *service.UserService
}

// NewUserHandler creates and returns a new UserHandler.
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

// HandleUsers manages requests for /users endpoint.
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

// HandleUserByID manages requests for /users/{id} endpoint.
func (h *UserHandler) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetUserByID(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetUsers handles GET requests for all users.
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := h.UserService.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUser handles POST requests to create a new user.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
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

// GetUserByID handles GET requests for a single user by ID.
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	user, err := h.UserService.GetUserByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
}
