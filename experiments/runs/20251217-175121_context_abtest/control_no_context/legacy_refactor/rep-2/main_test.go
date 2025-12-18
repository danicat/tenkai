package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"experiment/internal/api"
	"experiment/internal/service" // Import service to access User struct if needed for assertions
)

func TestUserAPI(t *testing.T) {
	userHandlers := &api.UserHandlers{}

	// Test POST /users
	t.Run("CreateUser", func(t *testing.T) {
		newUser := service.User{Name: "John Doe", Email: "john@example.com"}
		body, _ := json.Marshal(newUser)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		userHandlers.HandleUsers(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
		}

		var createdUser service.User
		json.NewDecoder(rec.Body).Decode(&createdUser)
		if createdUser.ID == "" {
			t.Error("Expected user ID to be generated, got empty")
		}
		if createdUser.Name != newUser.Name || createdUser.Email != newUser.Email {
			t.Errorf("Expected user %v, got %v", newUser, createdUser)
		}
	})

	// Test GET /users
	t.Run("GetAllUsers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()

		userHandlers.HandleUsers(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var users []service.User
		json.NewDecoder(rec.Body).Decode(&users)
		if len(users) != 1 { // Should be 1 user from the previous test
			t.Errorf("Expected 1 user, got %d", len(users))
		}
	})

	// Test GET /users/{id}
	t.Run("GetUserByID", func(t *testing.T) {
		// First, create a user to get an ID
		newUser := service.User{Name: "Jane Doe", Email: "jane@example.com"}
		body, _ := json.Marshal(newUser)
		postReq := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		postRec := httptest.NewRecorder()
		userHandlers.HandleUsers(postRec, postReq)
		var createdUser service.User
		json.NewDecoder(postRec.Body).Decode(&createdUser)
		userID := createdUser.ID

		req := httptest.NewRequest(http.MethodGet, "/users/"+userID, nil)
		rec := httptest.NewRecorder()

		userHandlers.HandleUserByID(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var retrievedUser service.User
		json.NewDecoder(rec.Body).Decode(&retrievedUser)
		if retrievedUser.ID != userID || retrievedUser.Name != newUser.Name {
			t.Errorf("Expected user %v, got %v", createdUser, retrievedUser)
		}
	})

	// Test GET /users/{id} not found
	t.Run("GetUserByIDNotFound", func(t *testing.T) {
		nonExistentID := "non-existent-id"
		req := httptest.NewRequest(http.MethodGet, "/users/"+nonExistentID, nil)
		rec := httptest.NewRecorder()

		userHandlers.HandleUserByID(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}
