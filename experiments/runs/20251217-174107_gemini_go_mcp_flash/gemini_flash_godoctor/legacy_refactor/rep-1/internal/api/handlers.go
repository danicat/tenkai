package api

import (
	"encoding/json"
	"net/http"

	"experiment/internal/service"
)

type UserHandlers struct {
	service *service.UserService
}

func NewUserHandlers(s *service.UserService) *UserHandlers {
	return &UserHandlers{service: s}
}

func (h *UserHandlers) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllUsers(w, r)
	case http.MethodPost:
		h.CreateUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandlers) HandleUserByID(w http.ResponseWriter, r *http.Request) {
id := r.URL.Path[len("/users/"):]

	switch r.Method {
	case http.MethodGet:
		h.GetUserByID(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.service.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandlers) GetUserByID(w http.ResponseWriter, r *http.Request, id string) {
	u, ok := h.service.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}
