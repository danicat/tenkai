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

func TestUsersAPI(t *testing.T) {
	// Initialize the user service
	userService := service.NewUserService()

	// Initialize API handlers with the user service
	handlers := api.NewUserHandlers(userService)

	// Test GET /users (empty initially)
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()
	handlers.UsersHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "[]\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Test POST /users
	user := service.User{Name: "Alice", Email: "alice@example.com"}
	jsonBody, _ := json.Marshal(user)
	req, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	rr = httptest.NewRecorder()
	handlers.UsersHandler(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code for POST: got %v want %v",
			status, http.StatusCreated)
	}

	var createdUser service.User
	json.NewDecoder(rr.Body).Decode(&createdUser)
	if createdUser.Name != "Alice" || createdUser.Email != "alice@example.com" || createdUser.ID == "" {
		t.Errorf("handler returned unexpected user for POST: %+v", createdUser)
	}

	// Test GET /users after POST
	req, _ = http.NewRequest(http.MethodGet, "/users", nil)
	rr = httptest.NewRecorder()
	handlers.UsersHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code after POST: got %v want %v",
			status, http.StatusOK)
	}

	var userList []service.User
	json.NewDecoder(rr.Body).Decode(&userList)
	if len(userList) != 1 || userList[0].ID != createdUser.ID {
		t.Errorf("handler returned unexpected user list after POST: %+v", userList)
	}

	// Test GET /users/{id}
	req, _ = http.NewRequest(http.MethodGet, "/users/"+createdUser.ID, nil)
	rr = httptest.NewRecorder()
	handlers.UserByIDHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for GET by ID: got %v want %v",
			status, http.StatusOK)
	}

	var fetchedUser service.User
	json.NewDecoder(rr.Body).Decode(&fetchedUser)
	if fetchedUser.ID != createdUser.ID || fetchedUser.Name != createdUser.Name {
		t.Errorf("handler returned unexpected user for GET by ID: %+v", fetchedUser)
	}

	// Test GET /users/{id} (not found)
	req, _ = http.NewRequest(http.MethodGet, "/users/nonexistent", nil)
	rr = httptest.NewRecorder()
	handlers.UserByIDHandler(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for GET by non-existent ID: got %v want %v",
			status, http.StatusNotFound)
	}
}
