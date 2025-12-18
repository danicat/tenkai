package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"experiment/internal/api"
	"experiment/internal/service"
)

func TestUsersHandler_GetUsers(t *testing.T) {
	// Clear existing users to ensure test isolation
	service.ClearUsers()

	// Add a test user
	service.CreateUser(service.User{Name: "Test User", Email: "test@example.com"})

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.UsersHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"id":"`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	var users []service.User
	if err := json.Unmarshal(rr.Body.Bytes(), &users); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}
	if len(users) == 0 || users[0].Name != "Test User" {
		t.Errorf("Expected at least one user named 'Test User', got %+v", users)
	}
}

func TestUsersHandler_CreateUser(t *testing.T) {
	// Clear existing users to ensure test isolation
	service.ClearUsers()

	var jsonStr = []byte(`{"Name":"New User","Email":"new@example.com"}`)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.UsersHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected := `{"id":"`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want substring %v",
			rr.Body.String(), expected)
	}

	var createdUser service.User
	if err := json.Unmarshal(rr.Body.Bytes(), &createdUser); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}
	if createdUser.Name != "New User" {
		t.Errorf("Expected created user name 'New User', got %s", createdUser.Name)
	}
}

func TestSpecificUserHandler_GetUserByID(t *testing.T) {
	// Clear existing users to ensure test isolation
	service.ClearUsers()

	// Add a test user
	userToFind, _ := service.CreateUser(service.User{Name: "Find Me", Email: "find@example.com"})

	req, err := http.NewRequest("GET", "/users/"+userToFind.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.SpecificUserHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var foundUser service.User
	if err := json.Unmarshal(rr.Body.Bytes(), &foundUser); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}
	if foundUser.ID != userToFind.ID || foundUser.Name != userToFind.Name {
		t.Errorf("Expected to find user %+v, got %+v", userToFind, foundUser)
	}
}

func TestSpecificUserHandler_NotFound(t *testing.T) {
	// Clear existing users to ensure test isolation
	service.ClearUsers()

	req, err := http.NewRequest("GET", "/users/nonexistent-id", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.SpecificUserHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}
