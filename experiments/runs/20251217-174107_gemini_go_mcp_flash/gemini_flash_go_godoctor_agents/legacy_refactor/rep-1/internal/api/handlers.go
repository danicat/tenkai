package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"experiment/internal/service"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	Service service.UserService
}

// NewUserHandler creates and returns a new UserHandler
func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

// ServeHTTP implements the http.Handler interface for /users endpoint
func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/users/") {
		h.handleGetUserByID(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetUsers(w, r)
	case http.MethodPost:
		h.handleCreateUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	list := h.Service.GetAllUsers()
	json.NewEncoder(w).Encode(list)
}

func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user service.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.Service.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (h *UserHandler) handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		http.NotFound(w, r) // No ID provided, so it's not a specific user request
		return
	}

	user, err := h.Service.GetUserByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
}
