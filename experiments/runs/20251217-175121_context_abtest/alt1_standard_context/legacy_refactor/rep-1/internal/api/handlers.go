package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"legacy_refactor/internal/service"
)

// UserHandlers provides HTTP handlers for user-related operations
type UserHandlers struct {
	service service.UserService
}

// NewUserHandlers creates a new instance of UserHandlers
func NewUserHandlers(s service.UserService) *UserHandlers {
	return &UserHandlers{service: s}
}

// RegisterRoutes registers the user API routes
func (h *UserHandlers) RegisterRoutes() {
	http.HandleFunc("/users", h.usersHandler)
	http.HandleFunc("/users/", h.userByIDHandler)
}

func (h *UserHandlers) usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandlers) userByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getUserByID(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandlers) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.service.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandlers) createUser(w http.ResponseWriter, r *http.Request) {
	var user service.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.service.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (h *UserHandlers) getUserByID(w http.ResponseWriter, r *http.Request, id string) {
	user, err := h.service.GetUserByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
}
