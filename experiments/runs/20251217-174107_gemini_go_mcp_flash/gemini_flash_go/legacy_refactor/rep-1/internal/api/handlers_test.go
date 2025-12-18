package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"experiment/internal/service"
)

func TestHandleUsers(t *testing.T) {
	svc := service.NewService()
	apiHandlers := NewAPI(svc)

	// Test POST /users (Create User)
	user := service.User{Name: "Alice", Email: "alice@example.com"}
	jsonUser, _ := json.Marshal(user)
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonUser))
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()

	apiHandlers.HandleUsers(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var createdUser service.User
	if err := json.NewDecoder(rec.Body).Decode(&createdUser); err != nil {
		t.Fatal(err)
	}
	if createdUser.ID == "" || createdUser.Name != "Alice" {
		t.Errorf("Unexpected created user: %v", createdUser)
	}

	// Test GET /users (Get All Users)
	req, err = http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec = httptest.NewRecorder()

	apiHandlers.HandleUsers(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var users []service.User
	if err := json.NewDecoder(rec.Body).Decode(&users); err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 || users[0].Name != "Alice" {
		t.Errorf("Unexpected users list: %v", users)
	}
}

func TestHandleUserByID(t *testing.T) {
	svc := service.NewService()
	apiHandlers := NewAPI(svc)

	// Create a user first
	user := service.User{Name: "Bob", Email: "bob@example.com"}
	createdUser := svc.CreateUser(user)

	// Test GET /users/{id} (Existing User)
	req, err := http.NewRequest(http.MethodGet, "/users/"+createdUser.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()

	apiHandlers.HandleUserByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var fetchedUser service.User
	if err := json.NewDecoder(rec.Body).Decode(&fetchedUser); err != nil {
		t.Fatal(err)
	}
	if fetchedUser.ID != createdUser.ID || fetchedUser.Name != createdUser.Name {
		t.Errorf("Unexpected fetched user: %v", fetchedUser)
	}

	// Test GET /users/{id} (Non-Existing User)
	req, err = http.NewRequest(http.MethodGet, "/users/999", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec = httptest.NewRecorder()

	apiHandlers.HandleUserByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
