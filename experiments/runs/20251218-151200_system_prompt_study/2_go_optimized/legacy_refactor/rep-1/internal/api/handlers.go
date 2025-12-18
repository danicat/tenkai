package api

import (
	"encoding/json"
	"net/http"

	"rep-1/internal/service"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userService service.UserServicer
}

// NewUserHandler creates and returns a new UserHandler
func NewUserHandler(s service.UserServicer) *UserHandler {
	return &UserHandler{
		userService: s,
	}
}

// HandleUsers manages requests to /users
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

// HandleUserByID manages requests to /users/{id}
func (h *UserHandler) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	switch r.Method {
	case http.MethodGet:
		h.GetUser(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetUsers handles GET requests for all users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := h.userService.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUser handles POST requests to create a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user service.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.userService.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// GetUser handles GET requests for a single user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, id string) {
	user, ok := h.userService.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
}
