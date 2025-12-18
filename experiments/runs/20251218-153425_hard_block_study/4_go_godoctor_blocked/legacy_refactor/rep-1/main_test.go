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

	// Test POST /users
	newUser := service.User{Name: "Alice", Email: "alice@example.com"}
	body, _ := json.Marshal(newUser)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	userHandlers.HandleUsers(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createdUser service.User
	json.NewDecoder(rr.Body).Decode(&createdUser)
	if createdUser.Name != "Alice" {
		t.Errorf("handler returned unexpected user name: got %v want %v",
			createdUser.Name, "Alice")
	}
	if createdUser.ID == "" {
		t.Errorf("handler returned user with empty ID")
	}

	// Test GET /users
	req, _ = http.NewRequest(http.MethodGet, "/users", nil)
	rr = httptest.NewRecorder()
	userHandlers.HandleUsers(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var users []service.User
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 1 {
		t.Errorf("handler returned unexpected number of users: got %v want %v",
			len(users), 1)
	}
	if users[0].Name != "Alice" {
		t.Errorf("handler returned unexpected user name: got %v want %v",
			users[0].Name, "Alice")
	}

	// Test GET /users/{id}
	req, _ = http.NewRequest(http.MethodGet, "/users/"+createdUser.ID, nil)
	rr = httptest.NewRecorder()
	userHandlers.HandleUserByID(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var fetchedUser service.User
	json.NewDecoder(rr.Body).Decode(&fetchedUser)
	if fetchedUser.ID != createdUser.ID {
		t.Errorf("handler returned unexpected user ID: got %v want %v",
			fetchedUser.ID, createdUser.ID)
	}

	// Test GET /users/{non-existent-id}
	req, _ = http.NewRequest(http.MethodGet, "/users/999", nil)
	rr = httptest.NewRecorder()
	userHandlers.HandleUserByID(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent user: got %v want %v",
			status, http.StatusNotFound)
	}
}
