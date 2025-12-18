package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"experiment/internal/service"
)

// API holds the service dependency and handles HTTP requests
type API struct {
	Service *service.Service
}

// NewAPI creates and returns a new API instance
func NewAPI(s *service.Service) *API {
	return &API{Service: s}
}

// HandleUsers handles requests to /users
func (a *API) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users := a.Service.GetAllUsers()
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var u service.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newUser, err := a.Service.CreateUser(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleUserByID handles requests to /users/{id}
func (a *API) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		u, ok := a.Service.GetUserByID(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(u)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
