package api

import (
	"encoding/json"
	"net/http"

	"rep-3/internal/service"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	Service service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

// ListUsers handles GET /users requests
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users := h.Service.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUser handles POST /users requests
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.Service.CreateUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// GetUser handles GET /users/{id} requests
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]

	u, ok := h.Service.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}

// SetupRoutes registers the user API routes
func (h *UserHandler) SetupRoutes() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.ListUsers(w, r)
		case http.MethodPost:
			h.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
