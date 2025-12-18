package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"experiment/internal/api"
	"experiment/internal/service" // Import service to reset state for testing
)

func TestUserAPI(t *testing.T) {
	// Reset the service state for each test run (important for consistent tests)
	service.ResetUsers()

	// Test POST /users
	user1 := service.User{Name: "Alice", Email: "alice@example.com"}
	user1JSON, _ := json.Marshal(user1)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(user1JSON))
	rr := httptest.NewRecorder()
	api.HandleUsers(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createdUser service.User
	json.NewDecoder(rr.Body).Decode(&createdUser)
	if createdUser.ID == "" || createdUser.Name != "Alice" {
		t.Errorf("handler returned unexpected body: %v", rr.Body.String())
	}

	// Test GET /users
	req, _ = http.NewRequest(http.MethodGet, "/users", nil)
	rr = httptest.NewRecorder()
	api.HandleUsers(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var users []service.User
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 1 || users[0].Name != "Alice" {
		t.Errorf("handler returned unexpected body for GET /users: %v", rr.Body.String())
	}

	// Test GET /users/{id}
	req, _ = http.NewRequest(http.MethodGet, "/users/"+createdUser.ID, nil)
	rr = httptest.NewRecorder()
	api.HandleUserByID(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var fetchedUser service.User
	json.NewDecoder(rr.Body).Decode(&fetchedUser)
	if fetchedUser.ID != createdUser.ID || fetchedUser.Name != "Alice" {
		t.Errorf("handler returned unexpected body for GET /users/{id}: %v", rr.Body.String())
	}

	// Test GET /users/{non-existent-id}
	req, _ = http.NewRequest(http.MethodGet, "/users/nonexistent", nil)
	rr = httptest.NewRecorder()
	api.HandleUserByID(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent user: got %v want %v",
			status, http.StatusNotFound)
	}
}
