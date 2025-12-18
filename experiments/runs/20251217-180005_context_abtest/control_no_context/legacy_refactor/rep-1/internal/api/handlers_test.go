package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rep-1/internal/service"
	"testing"

	"github.com/google/uuid"
)

func TestHandleUsers_Get(t *testing.T) {
	// Clear users map and add some test data for isolation
	service.ClearUsers()
	service.CreateUser("Alice", "alice@example.com")
	service.CreateUser("Bob", "bob@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	HandleUsers(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var users []service.User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	// Check if the names are present (order might vary)
	foundAlice := false
	foundBob := false
	for _, u := range users {
		if u.Name == "Alice" {
			foundAlice = true
		}
		if u.Name == "Bob" {
			foundBob = true
		}
	}

	if !foundAlice || !foundBob {
		t.Errorf("Expected to find Alice and Bob in the response")
	}
}

func TestHandleUsers_Post(t *testing.T) {
	// Clear users map before test
	service.ClearUsers()

	userPayload := []byte(`{"name": "Charlie", "email": "charlie@example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(userPayload))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	HandleUsers(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdUser service.User
	if err := json.NewDecoder(rr.Body).Decode(&createdUser); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if createdUser.Name != "Charlie" || createdUser.Email != "charlie@example.com" {
		t.Errorf("Expected created user to be Charlie, got %+v", createdUser)
	}
	if createdUser.ID == uuid.Nil {
		t.Errorf("Expected user ID to be non-nil")
	}

	// Verify user is actually in the service storage
	_, ok := service.GetUserByID(createdUser.ID)
	if !ok {
		t.Errorf("Expected user with ID %v to be found in service", createdUser.ID)
	}
}

func TestHandleUserByID_Get(t *testing.T) {
	// Clear users map and add a test user
	service.ClearUsers()
	user, _ := service.CreateUser("David", "david@example.com")

	req, _ := http.NewRequest(http.MethodGet, "/users/"+user.ID.String(), nil)
	rr := httptest.NewRecorder()

	HandleUserByID(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var retrievedUser service.User
	if err := json.NewDecoder(rr.Body).Decode(&retrievedUser); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if retrievedUser.ID != user.ID || retrievedUser.Name != user.Name {
		t.Errorf("Expected user %v, got %v", user, retrievedUser)
	}

	// Test case: User not found
	nonExistentID := uuid.New()
	reqNotFound, _ := http.NewRequest(http.MethodGet, "/users/"+nonExistentID.String(), nil)
	rrNotFound := httptest.NewRecorder()
	HandleUserByID(rrNotFound, reqNotFound)

	if status := rrNotFound.Code; status != http.StatusNotFound {
		t.Errorf("handler for non-existent user returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	// Test case: Invalid UUID
	reqInvalidID, _ := http.NewRequest(http.MethodGet, "/users/invalid-uuid", nil)
	rrInvalidID := httptest.NewRecorder()
	HandleUserByID(rrInvalidID, reqInvalidID)

	if status := rrInvalidID.Code; status != http.StatusBadRequest {
		t.Errorf("handler for invalid UUID returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
