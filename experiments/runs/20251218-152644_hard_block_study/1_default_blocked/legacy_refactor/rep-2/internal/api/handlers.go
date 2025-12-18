package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"experiment/internal/service"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// RegisterRoutes registers the user API routes
func (h *UserHandler) RegisterRoutes() {
	http.HandleFunc("/users", h.usersHandler)
	http.HandleFunc("/users/", h.userByIDHandler)
}

func (h *UserHandler) usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users := h.service.GetAllUsers()
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var u service.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u = h.service.CreateUser(u)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) userByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	switch r.Method {
	case http.MethodGet:
		u, ok := h.service.GetUserByID(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(u)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
