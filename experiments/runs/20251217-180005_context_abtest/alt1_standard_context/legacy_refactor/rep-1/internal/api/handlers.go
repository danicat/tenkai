package api

import (
	"encoding/json"
	"net/http"

	"legacy_refactor/internal/service"
)

// RegisterRoutes registers all HTTP handlers
func RegisterRoutes() {
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userByIDHandler)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list := service.GetAllUsers()
		json.NewEncoder(w).Encode(list)

	case http.MethodPost:
		var reqUser service.User
		if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		createdUser := service.CreateUser(reqUser.Name, reqUser.Email)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdUser)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func userByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	switch r.Method {
	case http.MethodGet:
		u, ok := service.GetUser(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(u)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
