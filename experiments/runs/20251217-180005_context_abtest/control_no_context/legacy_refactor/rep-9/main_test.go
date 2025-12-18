package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"experiment/internal/api"
	"experiment/internal/service"
)

func TestHandleUsers(t *testing.T) {
	// Clear existing users for a clean test environment
	service.ClearUsers() // Assuming a ClearUsers method exists or can be added

	router := http.NewServeMux()
	router.HandleFunc("/users", api.HandleUsers)

	// Test POST /users
	user := service.User{Name: "John Doe", Email: "john@example.com"}
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createdUser service.User
	json.NewDecoder(rr.Body).Decode(&createdUser)
	if createdUser.Name != user.Name || createdUser.Email != user.Email || createdUser.ID == "" {
		t.Errorf("handler returned unexpected body: got %+v expected name %s, email %s and non-empty ID", createdUser, user.Name, user.Email)
	}

	// Test GET /users
	req, _ = http.NewRequest("GET", "/users", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var users []service.User
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 1 || users[0].ID != createdUser.ID {
		t.Errorf("handler returned unexpected body: got %+v expected 1 user with ID %s", users, createdUser.ID)
	}
}

func TestHandleUserByID(t *testing.T) {
	// Clear existing users for a clean test environment
	service.ClearUsers() // Assuming a ClearUsers method exists or can be added

	router := http.NewServeMux()
	router.HandleFunc("/users/", api.HandleUserByID)

	// First, create a user to retrieve
	createdUser := service.CreateUser("Jane Doe", "jane@example.com")

	// Test GET /users/{id}
	req, _ := http.NewRequest("GET", "/users/"+createdUser.ID, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var retrievedUser service.User
	json.NewDecoder(rr.Body).Decode(&retrievedUser)
	if retrievedUser.ID != createdUser.ID || retrievedUser.Name != createdUser.Name {
		t.Errorf("handler returned unexpected body: got %+v expected %+v", retrievedUser, createdUser)
	}

	// Test GET /users/{non-existent-id}
	req, _ = http.NewRequest("GET", "/users/nonexistent", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent user: got %v want %v",
			status, http.StatusNotFound)
	}
}
