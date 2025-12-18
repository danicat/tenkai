package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"experiment/internal/service"
)

// Handlers struct holds the service dependency
type Handlers struct {
	Service *service.Service
}

// NewHandlers creates and initializes new Handlers
func NewHandlers(s *service.Service) *Handlers {
	return &Handlers{Service: s}
}

// HandleUsers handles requests to /users
func (h *Handlers) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list := h.Service.ListUsers()
		json.NewEncoder(w).Encode(list)

	case http.MethodPost:
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

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleUserByID handles requests to /users/{id}
func (h *Handlers) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	switch r.Method {
	case http.MethodGet:
		u, ok := h.Service.GetUser(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(u)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
