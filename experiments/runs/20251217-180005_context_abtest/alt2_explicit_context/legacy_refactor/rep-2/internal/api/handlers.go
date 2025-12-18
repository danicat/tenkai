package api

import (
	"encoding/json"
	"net/http"
	"strings"

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
		users := h.service.GetAllUsers()
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var u service.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u, _ = h.service.CreateUser(u)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandlers) HandleUserByID(w http.ResponseWriter, r *http.Request) {
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
