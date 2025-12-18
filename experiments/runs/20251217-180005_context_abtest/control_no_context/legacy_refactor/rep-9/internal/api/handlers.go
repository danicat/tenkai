package api

import (
	"encoding/json"
	"net/http"

	"experiment/internal/service"
)

// HandleUsers handles requests to /users
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list := service.ListUsers()
		json.NewEncoder(w).Encode(list)

	case http.MethodPost:
		var reqUser service.User
		if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u := service.CreateUser(reqUser.Name, reqUser.Email)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleUserByID handles requests to /users/{id}
func HandleUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	switch r.Method {
	case http.MethodGet:
		u, ok := service.GetUserByID(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(u)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
