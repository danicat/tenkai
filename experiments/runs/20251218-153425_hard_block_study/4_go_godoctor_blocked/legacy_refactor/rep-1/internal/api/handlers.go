package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"experiment/internal/service"
)

// UserHandlers provides http handlers for user-related requests
type UserHandlers struct {
	userService *service.UserService
}

// NewUserHandlers creates and returns a new UserHandlers
func NewUserHandlers(userService *service.UserService) *UserHandlers {
	return &UserHandlers{
		userService: userService,
	}
}

// HandleUsers manages requests to /users
func (h *UserHandlers) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users := h.userService.GetAllUsers()
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var user service.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		createdUser := h.userService.CreateUser(user)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdUser)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleUserByID manages requests to /users/{id}
func (h *UserHandlers) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		user, ok := h.userService.GetUserByID(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(user)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
