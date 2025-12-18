package api

import (
	"encoding/json"
	"net/http"

	"experiment/internal/service"
)

// UserHandlers handles HTTP requests for user-related operations
type UserHandlers struct {
	userService service.UserService
}

// NewUserHandlers creates a new instance of UserHandlers
func NewUserHandlers(us service.UserService) *UserHandlers {
	return &UserHandlers{userService: us}
}

// GetUsersHandler handles GET requests to /users
func (uh *UserHandlers) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := uh.userService.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

// CreateUserHandler handles POST requests to /users
func (uh *UserHandlers) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := uh.userService.CreateUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// GetUserByIDHandler handles GET requests to /users/{id}
func (uh *UserHandlers) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	u, ok := uh.userService.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}
