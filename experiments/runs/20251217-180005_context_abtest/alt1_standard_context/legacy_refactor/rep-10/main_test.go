package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"legacy_refactor/rep-10/internal/api"
	"legacy_refactor/rep-10/internal/service"
)

func TestUserAPI(t *testing.T) {
	// Initialize service and handlers for testing
	userService := service.NewUserService()
	userHandlers := api.NewUserHandlers(userService)

	// Register routes only for the test server
	mux := http.NewServeMux()
	mux.HandleFunc("/users", userHandlers.HandleUsers)
	mux.HandleFunc("/users/", userHandlers.HandleUserByID)

	// Test server
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Test POST /users (Create User)
	user := service.User{Name: "John Doe", Email: "john.doe@example.com"}
	jsonUser, _ := json.Marshal(user)

	res, err := ts.Client().Post(ts.URL+"/users", "application/json", bytes.NewBuffer(jsonUser))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d for POST, got %d", http.StatusCreated, res.StatusCode)
	}

	var createdUser service.User
	json.NewDecoder(res.Body).Decode(&createdUser)
	if createdUser.Name != user.Name || createdUser.Email != user.Email || createdUser.ID == "" {
		t.Errorf("Unexpected created user: %+v", createdUser)
	}

	// Test GET /users (Get All Users)
	res, err = ts.Client().Get(ts.URL + "/users")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d for GET, got %d", http.StatusOK, res.StatusCode)
	}

	var users []service.User
	json.NewDecoder(res.Body).Decode(&users)
	if len(users) != 1 || users[0].ID != createdUser.ID {
		t.Errorf("Unexpected users list: %+v", users)
	}

	// Test GET /users/{id} (Get User By ID)
	res, err = ts.Client().Get(ts.URL + "/users/" + createdUser.ID)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d for GET by ID, got %d", http.StatusOK, res.StatusCode)
	}

	var fetchedUser service.User
	json.NewDecoder(res.Body).Decode(&fetchedUser)
	if fetchedUser.ID != createdUser.ID || fetchedUser.Name != createdUser.Name {
		t.Errorf("Unexpected fetched user: %+v", fetchedUser)
	}

	// Test GET /users/{non-existent-id} (Get Non-Existent User)
	res, err = ts.Client().Get(ts.URL + "/users/non-existent-id")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status %d for non-existent user, got %d", http.StatusNotFound, res.StatusCode)
	}
}
