package api

import (
	"encoding/json"
	"net/http"

	"experiment/internal/service"
)

// UserHandlers handles HTTP requests related to users
type UserHandlers struct {
	service *service.UserService
}

// NewUserHandlers creates and returns new UserHandlers
func NewUserHandlers(s *service.UserService) *UserHandlers {
	return &UserHandlers{service: s}
}

// HandleUsers handles requests to /users
func (h *UserHandlers) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			h.GetAllUsers(w, r)
		case http.MethodPost:
			h.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetAllUsers handles GET /users
func (h *UserHandlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.service.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUser handles POST /users
func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
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

// HandleUserByID handles requests to /users/{id}
func (h *UserHandlers) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]

	switch r.Method {
		case http.MethodGet:
			h.GetUserByID(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetUserByID handles GET /users/{id}
func (h *UserHandlers) GetUserByID(w http.ResponseWriter, r *http.Request, id string) {
	user, ok := h.service.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
}
