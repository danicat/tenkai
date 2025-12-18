package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/petruzalek/rep-9/internal/service"
)

// Handler handles HTTP requests
type Handler struct {
	service *service.Service
}

// NewHandler creates and returns a new Handler instance
func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

// RegisterRoutes registers all HTTP handlers
func (h *Handler) RegisterRoutes() {
	http.HandleFunc("/users", h.usersHandler)
	http.HandleFunc("/users/", h.userByIDHandler)
}

func (h *Handler) usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) userByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getUser(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	list := h.service.ListUsers()
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.service.CreateUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request, id string) {
	user, ok := h.service.GetUser(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
}
