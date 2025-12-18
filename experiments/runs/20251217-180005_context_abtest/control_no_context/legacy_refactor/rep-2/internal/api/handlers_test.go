package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"experiment/internal/service"
)

func TestHandleUsers_Get(t *testing.T) {
	s := service.NewService()
	handlers := NewHandlers(s)

	user1 := service.User{Name: "User 1", Email: "user1@example.com"}
	user2 := service.User{Name: "User 2", Email: "user2@example.com"}
	s.CreateUser(user1)
	s.CreateUser(user2)

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandleUsers)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var users []service.User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

func TestHandleUsers_Post(t *testing.T) {
	s := service.NewService()
	handlers := NewHandlers(s)

	user := service.User{Name: "New User", Email: "new@example.com"}
	body, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/users", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandleUsers)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createdUser service.User
	if err := json.NewDecoder(rr.Body).Decode(&createdUser); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	if createdUser.Name != user.Name {
		t.Errorf("Expected created user name %s, got %s", user.Name, createdUser.Name)
	}
	if createdUser.ID == "" {
		t.Errorf("Expected ID to be generated, got empty")
	}
}

func TestHandleUserByID_Get(t *testing.T) {
	s := service.NewService()
	handlers := NewHandlers(s)

	user := service.User{Name: "Find Me", Email: "findme@example.com"}
	createdUser, _ := s.CreateUser(user)

	req, err := http.NewRequest("GET", "/users/"+createdUser.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandleUserByID)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var retrievedUser service.User
	if err := json.NewDecoder(rr.Body).Decode(&retrievedUser); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	if retrievedUser.ID != createdUser.ID {
		t.Errorf("Expected retrieved user ID %s, got %s", createdUser.ID, retrievedUser.ID)
	}
}

func TestHandleUserByID_NotFound(t *testing.T) {
	s := service.NewService()
	handlers := NewHandlers(s)

	req, err := http.NewRequest("GET", "/users/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandleUserByID)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}
