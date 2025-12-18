package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"rep-1/internal/service"
)

// HandleUsers handles requests for /users endpoint (GET and POST)
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users := service.GetAllUsers()
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var reqUser struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := service.CreateUser(reqUser.Name, reqUser.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleUserByID handles requests for /users/{id} endpoint (GET)
func HandleUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, ok := service.GetUserByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
}
