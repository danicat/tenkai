package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"rep-1/internal/api"
	"rep-1/internal/service"
)

func TestUserAPI(t *testing.T) {
	// Initialize service and handler layers for testing
	userService := service.NewUserService()
	userHandler := api.NewUserHandler(userService)

	// Create a new HTTP test server
	r := http.NewServeMux()
	r.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			userHandler.GetUsersHandler(w, req)
		case http.MethodPost:
			userHandler.CreateUserHandler(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	r.HandleFunc("/users/", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			userHandler.GetUserByIDHandler(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Test Case 1: Create a user
	newUser := service.User{Name: "Alice", Email: "alice@example.com"}
	body, _ := json.Marshal(newUser)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdUser service.User
	json.NewDecoder(rr.Body).Decode(&createdUser)
	if createdUser.Name != newUser.Name || createdUser.Email != newUser.Email || createdUser.ID == "" {
		t.Errorf("handler returned unexpected body: got %v want %v", createdUser, newUser)
	}
	
	// Test Case 2: Get all users
	req, _ = http.NewRequest(http.MethodGet, "/users", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var users []service.User
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 1 || users[0].ID != createdUser.ID {
		t.Errorf("handler returned unexpected body: got %v want a list with %v", users, createdUser)
	}

	// Test Case 3: Get user by ID
	req, _ = http.NewRequest(http.MethodGet, "/users/"+createdUser.ID, nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var fetchedUser service.User
	json.NewDecoder(rr.Body).Decode(&fetchedUser)
	if fetchedUser.ID != createdUser.ID {
		t.Errorf("handler returned unexpected body: got %v want %v", fetchedUser, createdUser)
	}

	// Test Case 4: Get non-existent user
	req, _ = http.NewRequest(http.MethodGet, "/users/nonexistent", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
