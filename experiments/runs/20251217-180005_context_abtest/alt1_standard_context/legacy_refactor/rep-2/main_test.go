package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"experiment/internal/api"
	"experiment/internal/service"
)

func TestUsersAPI(t *testing.T) {
	// Initialize service and handlers for testing
	svc := service.NewUserService()
	handlers := api.NewHandlers(svc)

	// Setup routes with our test handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/users", handlers.UsersHandler)
	mux.HandleFunc("/users/", handlers.UserByIDHandler)

	// Test Case 1: Create a user
	userJSON := `{"Name":"Alice", "Email":"alice@example.com"}`
	req, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); !strings.HasPrefix(contentType, expectedContentType) {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}

	var createdUser service.User
	if err := json.NewDecoder(rr.Body).Decode(&createdUser); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if createdUser.Name != "Alice" || createdUser.Email != "alice@example.com" || createdUser.ID == "" {
		t.Errorf("handler returned unexpected user: %+v", createdUser)
	}

	// Test Case 2: Get all users
	req, _ = http.NewRequest(http.MethodGet, "/users", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var users []service.User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if len(users) != 1 || users[0].ID != createdUser.ID {
		t.Errorf("handler returned unexpected users list: %+v", users)
	}

	// Test Case 3: Get user by ID
	req, _ = http.NewRequest(http.MethodGet, "/users/"+createdUser.ID, nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var retrievedUser service.User
	if err := json.NewDecoder(rr.Body).Decode(&retrievedUser); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if retrievedUser.ID != createdUser.ID || retrievedUser.Name != createdUser.Name || retrievedUser.Email != createdUser.Email {
		t.Errorf("handler returned unexpected user by ID: %+v", retrievedUser)
	}

	// Test Case 4: Get non-existent user by ID
	req, _ = http.NewRequest(http.MethodGet, "/users/nonexistent", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent user: got %v want %v", status, http.StatusNotFound)
	}
}
