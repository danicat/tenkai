package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"legacy_refactor/rep-1/internal/api"
	"legacy_refactor/rep-1/internal/service"
)

func TestUserAPI(t *testing.T) {
	svc := service.NewService()
	handler := api.NewHandler(svc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Test Create User
	user := service.User{Name: "John Doe", Email: "john@example.com"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var createdUser service.User
	json.NewDecoder(rr.Body).Decode(&createdUser)
	if createdUser.ID == "" {
		t.Error("Expected a non-empty user ID, got empty")
	}

	// Test List Users
	req = httptest.NewRequest(http.MethodGet, "/users", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var users []service.User
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

	// Test Get User by ID
	req = httptest.NewRequest(http.MethodGet, "/users/"+createdUser.ID, nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var fetchedUser service.User
	json.NewDecoder(rr.Body).Decode(&fetchedUser)
	if fetchedUser.ID != createdUser.ID {
		t.Errorf("Expected user ID %s, got %s", createdUser.ID, fetchedUser.ID)
	}
	if fetchedUser.Name != createdUser.Name {
		t.Errorf("Expected user name %s, got %s", createdUser.Name, fetchedUser.Name)
	}

	// Test Get Non-Existent User by ID
	req = httptest.NewRequest(http.MethodGet, "/users/non-existent-id", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for non-existent user, got %d", http.StatusNotFound, rr.Code)
	}
}
