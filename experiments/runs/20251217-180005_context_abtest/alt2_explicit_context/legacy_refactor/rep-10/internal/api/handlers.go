package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"experiment/internal/service"
)

type UserHandlers struct {
	UserService service.UserServicer
}

func (h *UserHandlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := h.UserService.GetUsers()
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.UserService.CreateUser(u.Name, u.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (h *UserHandlers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	u, err := h.UserService.GetUser(id)
	if err != nil {
		if err == service.ErrUserNotFound {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(u)
}
