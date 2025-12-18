package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"experiment/internal/service"
)

// UserHandlers holds the dependencies for user HTTP handlers
type UserHandlers struct {
	UserService *service.UserService
}

// NewUserHandlers creates a new UserHandlers instance
func NewUserHandlers(us *service.UserService) *UserHandlers {
	return &UserHandlers{UserService: us}
}

// HandleUsers manages requests for /users endpoint
func (uh *UserHandlers) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		uh.GetUsers(w, r)
	case http.MethodPost:
		uh.CreateUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleUserByID manages requests for /users/{id} endpoint
func (uh *UserHandlers) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	switch r.Method {
	case http.MethodGet:
		uh.GetUserByID(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetUsers handles GET /users
func (uh *UserHandlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := uh.UserService.GetUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUser handles POST /users
func (uh *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := uh.UserService.CreateUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// GetUserByID handles GET /users/{id}
func (uh *UserHandlers) GetUserByID(w http.ResponseWriter, r *http.Request, id string) {
	u, ok := uh.UserService.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}
