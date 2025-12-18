package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rep-2/internal/api"
	"rep-2/internal/service"
	"testing"
)

func TestUserAPI(t *testing.T) {
	// Initialize service and handlers for testing
	userService := service.NewUserService()
	handler := api.NewUserHandler(userService)
	handler.RegisterRoutes()

	// Test Case 1: Create a user (POST /users)
	user1 := service.User{Name: "Alice", Email: "alice@example.com"}
	user1JSON, _ := json.Marshal(user1)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(user1JSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdUser1 service.User
	if err := json.NewDecoder(rr.Body).Decode(&createdUser1); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if createdUser1.Name != user1.Name || createdUser1.Email != user1.Email || createdUser1.ID == "" {
		t.Errorf("handler returned unexpected body: got %v want ID, Name and Email match", createdUser1)
	}

	// Test Case 2: Get all users (GET /users)
	req, _ = http.NewRequest(http.MethodGet, "/users", nil)
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var users []service.User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(users) != 1 || users[0].ID != createdUser1.ID {
		t.Errorf("handler returned unexpected user list: got %v want list with one user %v", users, createdUser1)
	}

	// Test Case 3: Create another user
	user2 := service.User{Name: "Bob", Email: "bob@example.com"}
	user2JSON, _ := json.Marshal(user2)
	req, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(user2JSON))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	var createdUser2 service.User
	json.NewDecoder(rr.Body).Decode(&createdUser2)

	// Test Case 4: Get user by ID (GET /users/:id)
	req, _ = http.NewRequest(http.MethodGet, "/users/"+createdUser1.ID, nil)
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var fetchedUser service.User
	if err := json.NewDecoder(rr.Body).Decode(&fetchedUser); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if fetchedUser.ID != createdUser1.ID {
		t.Errorf("handler returned wrong user: got %v want %v", fetchedUser, createdUser1)
	}

	// Test Case 5: Get non-existent user by ID (GET /users/:id)
	req, _ = http.NewRequest(http.MethodGet, "/users/non-existent-id", nil)
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent user: got %v want %v", status, http.StatusNotFound)
	}
}
