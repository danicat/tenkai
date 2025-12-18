package api

import (
	"encoding/json"
	"net/http"

	"experiment/internal/service"
)

// User represents a user in the system
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// HandleUsers handles requests to /users
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list := service.GetAllUsers()
		json.NewEncoder(w).Encode(list)

	case http.MethodPost:
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		serviceUser := service.User{Name: u.Name, Email: u.Email}
		createdUser := service.CreateUser(serviceUser)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdUser)

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
