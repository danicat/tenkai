package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"experiment/internal/service"
)

// HandleUsers handles requests to /users
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list := service.ListUsers()
		json.NewEncoder(w).Encode(list)

	case http.MethodPost:
		var u service.User // Use service.User here
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		createdUser, err := service.CreateUser(u.Name, u.Email)
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

// HandleUserByID handles requests to /users/{id}
func HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		http.Error(w, "User ID not provided", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		u, err := service.GetUser(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(u)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
