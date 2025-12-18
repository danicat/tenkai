package api

import (
	"encoding/json"
	"net/http"

	"experiment/internal/service"
)

// UserHandler handles HTTP requests for user-related operations
type UserHandler struct {
	service service.UserService
}

// NewUserHandler creates and returns a new UserHandler
func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

// RegisterRoutes registers the user API routes
func (h *UserHandler) RegisterRoutes() {
	http.HandleFunc("/users", h.handleUsers)
	http.HandleFunc("/users/", h.handleUserByID)
}

func (h *UserHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllUsersHandler(w, r)
	case http.MethodPost:
		h.CreateUserHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) handleUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]

	switch r.Method {
	case http.MethodGet:
		h.GetUserByIDHandler(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetAllUsersHandler handles GET requests for all users
func (h *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := h.service.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUserHandler handles POST requests to create a new user
func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user service.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.service.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // Or appropriate error
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// GetUserByIDHandler handles GET requests for a single user by ID
func (h *UserHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request, id string) {
	user, err := h.service.GetUserByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
}
