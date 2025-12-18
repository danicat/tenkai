package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"rep-3/internal/service"
)

func TestGetUsers(t *testing.T) {
	// Clear users and add some test data
	service.ClearUsers()
	service.CreateUser("Alice", "alice@example.com")
	service.CreateUser("Bob", "bob@example.com")

	req, _ := http.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()

	GetUsers(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var users []service.User
	json.NewDecoder(rr.Body).Decode(&users)

	if len(users) != 2 {
		t.Errorf("handler returned unexpected number of users: got %v want %v", len(users), 2)
	}
}

func TestCreateUser(t *testing.T) {
	// Clear users before test
	service.ClearUsers()

	var jsonStr = []byte(`{"name":"Charlie","email":"charlie@example.com"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateUser(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var user service.User
	json.NewDecoder(rr.Body).Decode(&user)

	if user.Name != "Charlie" || user.Email != "charlie@example.com" {
		t.Errorf("handler returned unexpected user: %+v", user)
	}

	if _, ok := service.GetUserByID(user.ID); !ok {
		t.Errorf("User not found in service after creation via handler")
	}
}

func TestGetUser(t *testing.T) {
	// Clear users and add a test user
	service.ClearUsers()
	user := service.CreateUser("David", "david@example.com")

	req, _ := http.NewRequest("GET", "/users/"+user.ID, nil)
	rr := httptest.NewRecorder()

	GetUser(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var foundUser service.User
	json.NewDecoder(rr.Body).Decode(&foundUser)

	if foundUser.ID != user.ID || foundUser.Name != user.Name || foundUser.Email != user.Email {
		t.Errorf("Found user mismatch. Expected %+v, got %+v", user, foundUser)
	}

	// Test for non-existent user
	req, _ = http.NewRequest("GET", "/users/nonexistent-id", nil)
	rr = httptest.NewRecorder()
	GetUser(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for nonexistent user: got %v want %v", status, http.StatusNotFound)
	}
}
