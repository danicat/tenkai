package api

import (
	"encoding/json"
	"net/http"

	"experiment/internal/service"
)

// UserHandlers handles HTTP requests for user-related operations
type UserHandlers struct {
	service service.UserService
}

// NewUserHandlers creates a new instance of UserHandlers
func NewUserHandlers(s service.UserService) *UserHandlers {
	return &UserHandlers{service: s}
}

// RegisterRoutes registers the user-related HTTP routes
func (h *UserHandlers) RegisterRoutes() {
	http.HandleFunc("/users", h.usersHandler)
	http.HandleFunc("/users/", h.userByIDHandler)
}

func (h *UserHandlers) usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list := h.service.GetUsers()
		json.NewEncoder(w).Encode(list)

	case http.MethodPost:
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

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandlers) userByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]

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
