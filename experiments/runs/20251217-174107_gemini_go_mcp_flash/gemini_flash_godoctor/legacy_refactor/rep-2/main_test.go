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
	userService := service.NewUserService()
	userHandlers := api.NewUserHandlers(userService)

	rr := httptest.NewRecorder()

	// Test CreateUser
	newUser := service.User{Name: "Alice", Email: "alice@example.com"}
	body, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))

	userHandlers.CreateUserHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var createdUser service.User
	json.NewDecoder(rr.Body).Decode(&createdUser)

	if createdUser.Name != newUser.Name || createdUser.Email != newUser.Email || createdUser.ID == "" {
		t.Errorf("Created user mismatch: %+v", createdUser)
	}

	// Test GetUsers
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/users", nil)

	userHandlers.GetUsersHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var users []service.User
	json.NewDecoder(rr.Body).Decode(&users)

	if len(users) != 1 || users[0].ID != createdUser.ID {
		t.Errorf("Get all users mismatch: %+v", users)
	}

	// Test GetUserByID
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/users/"+createdUser.ID, nil)

	userHandlers.GetUserByIDHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var fetchedUser service.User
	json.NewDecoder(rr.Body).Decode(&fetchedUser)

	if fetchedUser.ID != createdUser.ID {
		t.Errorf("Fetched user by ID mismatch: %+v", fetchedUser)
	}

	// Test GetUserByID (not found)
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/users/nonexistent", nil)

	userHandlers.GetUserByIDHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}
