package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"experiment/internal/service"
)

// UserHandlers provides HTTP handlers for user-related operations
type UserHandlers struct {
	Service service.UserService
}

// NewUserHandlers creates and returns new UserHandlers
func NewUserHandlers(s service.UserService) *UserHandlers {
	return &UserHandlers{Service: s}
}

// HandleUsers manages requests to /users
func (h *UserHandlers) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users := h.Service.GetUsers()
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var u service.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newUser, err := h.Service.CreateUser(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)

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
		u, err := h.Service.GetUserByID(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(u)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
