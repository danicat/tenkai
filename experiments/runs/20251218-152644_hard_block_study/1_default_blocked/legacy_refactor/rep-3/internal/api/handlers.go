package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"experiment/internal/service"
)

type UserHandlers struct {
	service service.UserService
}

func NewUserHandlers(s service.UserService) *UserHandlers {
	return &UserHandlers{service: s}
}

func (uh *UserHandlers) RegisterRoutes() {
	http.HandleFunc("/users", uh.usersHandler)
	http.HandleFunc("/users/", uh.userByIDHandler)
}

func (uh *UserHandlers) usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list := uh.service.GetUsers()
		json.NewEncoder(w).Encode(list)

	case http.MethodPost:
		var u service.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u, err := uh.service.CreateUser(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (uh *UserHandlers) userByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		u, ok := uh.service.GetUserByID(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(u)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
