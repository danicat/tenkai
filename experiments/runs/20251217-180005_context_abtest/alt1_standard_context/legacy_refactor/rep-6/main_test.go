package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"rep-6/internal/api"
	"rep-6/internal/service"
)

func TestUsersHandler(t *testing.T) {
	// Reset the service state for each test
	service.ResetUsersForTest() // Need to add this function to the service package

	// Test POST /users
	newUser := service.User{Name: "Alice", Email: "alice@example.com"}
	body, _ := json.Marshal(newUser)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	rw := httptest.NewRecorder()
	api.UsersHandler(rw, req)

	if rw.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rw.Code)
	}

	var createdUser service.User
	json.NewDecoder(rw.Body).Decode(&createdUser)
	if createdUser.Name != newUser.Name || createdUser.Email != newUser.Email || createdUser.ID == "" {
		t.Errorf("Created user mismatch: %+v", createdUser)
	}

	// Test GET /users
	req = httptest.NewRequest(http.MethodGet, "/users", nil)
	rw = httptest.NewRecorder()
	api.UsersHandler(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rw.Code)
	}

	var usersList []service.User
	json.NewDecoder(rw.Body).Decode(&usersList)
	if len(usersList) != 1 || usersList[0].ID != createdUser.ID {
		t.Errorf("Users list mismatch: %+v", usersList)
	}

	// Test GET /users/{id}
	req = httptest.NewRequest(http.MethodGet, "/users/"+createdUser.ID, nil)
	rw = httptest.NewRecorder()
	api.UserByIDHandler(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rw.Code)
	}

	var fetchedUser service.User
	json.NewDecoder(rw.Body).Decode(&fetchedUser)
	if fetchedUser.ID != createdUser.ID {
		t.Errorf("Fetched user by ID mismatch: %+v", fetchedUser)
	}
}
