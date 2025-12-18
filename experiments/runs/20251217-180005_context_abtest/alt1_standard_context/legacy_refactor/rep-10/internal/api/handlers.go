package api

import (
	"encoding/json"
	"net/http"

	"legacy_refactor/rep-10/internal/service"
)

// UserHandlers holds the dependencies for user HTTP handlers
type UserHandlers struct {
	UserService service.UserService
}

// NewUserHandlers creates and returns new UserHandlers
func NewUserHandlers(userService service.UserService) *UserHandlers {
	return &UserHandlers{
		UserService: userService,
	}
}

// HandleUsers manages requests to /users
func (h *UserHandlers) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetUsers(w, r)
	case http.MethodPost:
		h.CreateUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetUsers handles GET requests for /users
func (h *UserHandlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := h.UserService.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUser handles POST requests for /users
func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
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

// HandleUserByID manages requests to /users/{id}
func (h *UserHandlers) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	switch r.Method {
	case http.MethodGet:
		h.GetUserByID(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetUserByID handles GET requests for /users/{id}
func (h *UserHandlers) GetUserByID(w http.ResponseWriter, r *http.Request, id string) {
	user, ok := h.UserService.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
}
