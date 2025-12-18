package api

import (
	"encoding/json"
	"net/http"
	"experiment/internal/service"
)

// UserHandlers provides http handlers for user-related requests
type UserHandlers struct {
	// No explicit service field needed as service functions are package-level
}

func (h *UserHandlers) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.HandleGetUsers(w, r)
	case http.MethodPost:
		h.HandlePostUsers(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandlers) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	list := service.GetAllUsers()
	json.NewEncoder(w).Encode(list)
}

func (h *UserHandlers) HandlePostUsers(w http.ResponseWriter, r *http.Request) {
	var u service.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u = service.CreateUser(u)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandlers) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	switch r.Method {
	case http.MethodGet:
		h.HandleGetUserByID(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandlers) HandleGetUserByID(w http.ResponseWriter, r *http.Request, id string) {
	u, ok := service.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(u)
}
