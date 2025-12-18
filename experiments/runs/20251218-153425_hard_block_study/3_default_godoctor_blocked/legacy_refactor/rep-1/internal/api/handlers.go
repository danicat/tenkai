package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"experiment/internal/service"
)

// UserService defines the interface for user-related business logic.
type UserService interface {
	GetAllUsers() []service.User
	CreateUser(user service.User) (service.User, error)
	GetUserByID(id string) (service.User, bool)
}

type UserHandlers struct {
	userService UserService
}

func NewUserHandlers(userService UserService) *UserHandlers {
	return &UserHandlers{userService: userService}
}

func (h *UserHandlers) UsersHandler(w http.ResponseWriter, r *http.Request) {
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

		createdUser, err := h.userService.CreateUser(user)
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

func (h *UserHandlers) UserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, ok := h.userService.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
}
