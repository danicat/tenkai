package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"experiment/internal/api"
	"experiment/internal/service"
)

func TestGetUsersHandler(t *testing.T) {
	service.ClearUsers()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users", nil)

	api.GetUsersHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var users []service.User
	if err := json.NewDecoder(w.Body).Decode(&users); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}

	service.CreateUser(service.User{Name: "Test User", Email: "test@example.com"})

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/users", nil)

	api.GetUsersHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if err := json.NewDecoder(w.Body).Decode(&users); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}

func TestCreateUserHandler(t *testing.T) {
	service.ClearUsers()

	u := service.User{Name: "New User", Email: "new@example.com"}
	body, _ := json.Marshal(u)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))

	api.CreateUserHandler(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var createdUser service.User
	if err := json.NewDecoder(w.Body).Decode(&createdUser); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createdUser.ID == "" {
		t.Error("expected user ID to be generated")
	}
	if createdUser.Name != u.Name {
		t.Errorf("expected user name %q, got %q", u.Name, createdUser.Name)
	}
}

func TestGetUserByIDHandler(t *testing.T) {
	service.ClearUsers()

	u := service.User{Name: "Find Me", Email: "findme@example.com"}
	createdUser, _ := service.CreateUser(u)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/"+createdUser.ID, nil)

	api.GetUserByIDHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var foundUser service.User
	if err := json.NewDecoder(w.Body).Decode(&foundUser); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if foundUser.ID != createdUser.ID {
		t.Errorf("expected user ID %q, got %q", createdUser.ID, foundUser.ID)
	}

	// Test non-existent user
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/users/nonexistent-id", nil)

	api.GetUserByIDHandler(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d for non-existent user, got %d", http.StatusNotFound, w.Code)
	}
}
