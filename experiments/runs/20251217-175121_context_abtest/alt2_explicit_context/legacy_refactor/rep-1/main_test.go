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
	// Initialize service and handlers for testing
	svc := service.NewUserService()
	handlers := api.NewUserHandlers(svc)
	handlers.RegisterRoutes()

	// Test case 1: Create a user (POST /users)
	newUser := service.User{Name: "Alice", Email: "alice@example.com"}
	body, _ := json.Marshal(newUser)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createdUser service.User
	json.NewDecoder(rr.Body).Decode(&createdUser)
	if createdUser.Name != newUser.Name || createdUser.Email != newUser.Email || createdUser.ID == "" {
		t.Errorf("handler returned unexpected body: got %v want %v with an ID",
			rr.Body.String(), newUser)
	}

	// Test case 2: Get all users (GET /users)
	req, _ = http.NewRequest(http.MethodGet, "/users", nil)
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var users []service.User
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 1 || users[0].ID != createdUser.ID {
		t.Errorf("handler returned unexpected body for GET /users: got %v, expected a list with %v",
			users, createdUser)
	}

	// Test case 3: Get user by ID (GET /users/{id})
	req, _ = http.NewRequest(http.MethodGet, "/users/"+createdUser.ID, nil)
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var fetchedUser service.User
	json.NewDecoder(rr.Body).Decode(&fetchedUser)
	if fetchedUser.ID != createdUser.ID || fetchedUser.Name != createdUser.Name {
		t.Errorf("handler returned unexpected body for GET /users/{id}: got %v, expected %v",
			fetchedUser, createdUser)
	}

	// Test case 4: Get non-existent user by ID (GET /users/{id})
	req, _ = http.NewRequest(http.MethodGet, "/users/non-existent-id", nil)
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent user: got %v want %v",
			status, http.StatusNotFound)
	}
}
