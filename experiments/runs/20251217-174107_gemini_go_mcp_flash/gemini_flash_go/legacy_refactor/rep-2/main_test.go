package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"legacy_refactor/rep-2/internal/api"
	"legacy_refactor/rep-2/internal/service"
)

func TestUserAPI(t *testing.T) {
	// Initialize service and handlers
	userService := service.NewUserService()
	apiHandlers := api.NewHandlers(userService)

	// Setup HTTP server for testing
	mux := http.NewServeMux()
	mux.HandleFunc("/users", apiHandlers.HandleUsers)
	mux.HandleFunc("/users/", apiHandlers.HandleUserByID)

	// Test Create User (POST /users)
	user := service.User{Name: "John Doe", Email: "john@example.com"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rw.Code)
	}

	var createdUser service.User
	json.NewDecoder(rw.Body).Decode(&createdUser)
	if createdUser.ID == "" {
		t.Error("Expected user ID to be set, got empty")
	}

	// Test Get All Users (GET /users)
	req = httptest.NewRequest(http.MethodGet, "/users", nil)
	rw = httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rw.Code)
	}

	var users []service.User
	json.NewDecoder(rw.Body).Decode(&users)
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}
	if users[0].ID != createdUser.ID {
		t.Errorf("Expected user ID %s, got %s", createdUser.ID, users[0].ID)
	}

	// Test Get User By ID (GET /users/{id})
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", createdUser.ID), nil)
	rw = httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rw.Code)
	}

	var fetchedUser service.User
	json.NewDecoder(rw.Body).Decode(&fetchedUser)
	if fetchedUser.ID != createdUser.ID {
		t.Errorf("Expected user ID %s, got %s", createdUser.ID, fetchedUser.ID)
	}

	// Test Get Non-Existent User By ID (GET /users/{id})
	req = httptest.NewRequest(http.MethodGet, "/users/non-existent-id", nil)
	rw = httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d for non-existent user", http.StatusNotFound, rw.Code)
	}
}
