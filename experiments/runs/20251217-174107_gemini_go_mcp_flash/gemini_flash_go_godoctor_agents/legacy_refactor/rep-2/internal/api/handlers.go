package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"experiment/internal/service"
)

// Handler handles HTTP requests
type Handler struct {
	Service *service.Service
}

// NewHandler creates a new instance of the Handler
func NewHandler(s *service.Service) *Handler {
	return &Handler{Service: s}
}

// HandleUsers handles requests for /users
func (h *Handler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users := h.Service.GetAllUsers()
		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Printf("error encoding users: %v", err) // Log the error, but don't fail the request
		}

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
		if err := json.NewEncoder(w).Encode(createdUser); err != nil {
			log.Printf("error encoding created user: %v", err)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleUserByID handles requests for /users/{id}
func (h *Handler) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	switch r.Method {
	case http.MethodGet:
		u, err := h.Service.GetUserByID(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(u); err != nil {
			log.Printf("error encoding user by ID: %v", err)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
