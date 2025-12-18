package api

import (
	"encoding/json"
	"net/http"

	"rep-1/internal/service"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	UserService *service.UserService
}

// NewUserHandler creates and returns a new UserHandler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

// GetUsersHandler handles GET /users requests
func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	list := h.UserService.GetAllUsers()
	json.NewEncoder(w).Encode(list)
}

// CreateUserHandler handles POST /users requests
func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser := h.UserService.CreateUser(u)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// GetUserByIDHandler handles GET /users/{id} requests
func (h *UserHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	u, ok := h.UserService.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}