package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/petruzalek/rep-9/internal/api"
	"github.com/petruzalek/rep-9/internal/service"
)

func TestUsersAPI(t *testing.T) {
	svc := service.NewService()
	handler := api.NewHandler(svc)
	handler.RegisterRoutes()

	// Test POST /users (Create User)
	user := service.User{Name: "John Doe", Email: "john.doe@example.com"}
	jsonUser, _ := json.Marshal(user)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdUser service.User
	json.NewDecoder(rr.Body).Decode(&createdUser)
	if createdUser.Name != user.Name || createdUser.Email != user.Email || createdUser.ID == "" {
		t.Errorf("handler returned unexpected body: got %v expected %v with ID", createdUser, user)
	}
	fmt.Printf("Created user: %+v\n", createdUser.ID)

	// Test GET /users (List Users)
	req, _ = http.NewRequest(http.MethodGet, "/users", nil)
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var users []service.User
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 1 || users[0].ID != createdUser.ID {
		t.Errorf("handler returned unexpected body: got %v expected list with created user", users)
	}

	// Test GET /users/{id} (Get User by ID)
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", createdUser.ID), nil)
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var fetchedUser service.User
	json.NewDecoder(rr.Body).Decode(&fetchedUser)
	if fetchedUser.ID != createdUser.ID {
		t.Errorf("handler returned unexpected body: got %v expected %v", fetchedUser, createdUser)
	}

	// Test GET /users/{non-existent-id}
	req, _ = http.NewRequest(http.MethodGet, "/users/nonexistent", nil)
	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent user: got %v want %v", status, http.StatusNotFound)
	}
}
