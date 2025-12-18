package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"experiment/internal/service"
)

// API holds the dependencies for our HTTP handlers.
type API struct {
	service *service.Service
}

// NewAPI creates and returns a new API instance.
func NewAPI(s *service.Service) *API {
	return &API{service: s}
}

// HandleUsers manages requests to the /users endpoint.
func (a *API) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.getAllUsers(w, r)
	case http.MethodPost:
		a.createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *API) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users := a.service.GetAllUsers()
	json.NewEncoder(w).Encode(users)
}

func (a *API) createUser(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser := a.service.CreateUser(u)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// HandleUserByID manages requests to the /users/{id} endpoint.
func (a *API) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	switch r.Method {
	case http.MethodGet:
		a.getUserByID(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *API) getUserByID(w http.ResponseWriter, r *http.Request, id string) {
	u, ok := a.service.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}
