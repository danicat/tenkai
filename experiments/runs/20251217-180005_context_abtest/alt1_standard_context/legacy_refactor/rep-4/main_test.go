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

func TestUserAPI(t *testing.T) {
	// Initialize service and handlers
	userService := service.NewUserService()
	userHandlers := api.NewUserHandlers(userService)

	// Create a new HTTP test server
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)

	// Test GET /users (empty initially)
	userHandlers.HandleUsers(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var users []service.User
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(users))
	}

	// Test POST /users
	newUser := service.User{Name: "Alice", Email: "alice@example.com"}
	jsonBody, _ := json.Marshal(newUser)
	req, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	rr = httptest.NewRecorder()
	userHandlers.HandleUsers(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var createdUser service.User
	json.NewDecoder(rr.Body).Decode(&createdUser)
	if createdUser.Name != newUser.Name || createdUser.Email != newUser.Email || createdUser.ID == "" {
		t.Errorf("Created user mismatch: %+v", createdUser)
	}

	// Test GET /users again (with one user)
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/users", nil)
	userHandlers.HandleUsers(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	users = []service.User{}
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

	// Test GET /users/{id}
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/users/"+createdUser.ID, nil)
	userHandlers.HandleUserByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var fetchedUser service.User
	json.NewDecoder(rr.Body).Decode(&fetchedUser)
	if fetchedUser.ID != createdUser.ID || fetchedUser.Name != createdUser.Name {
		t.Errorf("Fetched user mismatch: %+v", fetchedUser)
	}

	// Test GET /users/{non-existent-id}
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/users/non-existent", nil)
	userHandlers.HandleUserByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}
