package api

import (
	"encoding/json"
	"net/http"

	"rep-6/internal/service"
)

// User represents a user in the system. This struct is duplicated from the service package.
// In a larger application, a dedicated `model` package would be preferred for shared structs.
// For this refactor, I will reuse the `service.User` struct directly in handlers to avoid duplication
// and keep it consistent with the original monolith's tight coupling.

// UsersHandler handles requests to /users
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list := service.GetAllUsers()
		json.NewEncoder(w).Encode(list)

	case http.MethodPost:
		var u service.User // Use service.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u = service.CreateUser(u)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// UserByIDHandler handles requests to /users/{id}
func UserByIDHandler(w http.ResponseWriter, r *http.Request) {
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
