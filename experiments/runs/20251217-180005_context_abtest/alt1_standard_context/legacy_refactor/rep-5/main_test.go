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
	api.RegisterRoutes()

	// Test POST /users
	newUser := service.User{Name: "John Doe", Email: "john@example.com"}
	body, _ := json.Marshal(newUser)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdUser service.User
	json.NewDecoder(rr.Body).Decode(&createdUser)
	if createdUser.Name != newUser.Name || createdUser.Email != newUser.Email || createdUser.ID == "" {
		t.Errorf("handler returned unexpected body: got %v want ID, Name, Email set", createdUser)
	}

	// Test GET /users
	req, _ = http.NewRequest(http.MethodGet, "/users", nil)
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var users []service.User
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 1 {
		t.Errorf("handler returned unexpected number of users: got %v want %v", len(users), 1)
	}

	// Test GET /users/{id}
	req, _ = http.NewRequest(http.MethodGet, "/users/"+createdUser.ID, nil)
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var fetchedUser service.User
	json.NewDecoder(rr.Body).Decode(&fetchedUser)
	if fetchedUser.ID != createdUser.ID {
		t.Errorf("handler returned unexpected user: got %v want %v", fetchedUser, createdUser)
	}
}
