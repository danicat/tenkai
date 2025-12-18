package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"rep-5/internal/api"
	"rep-5/internal/service"
)

func TestGetUsersHandler(t *testing.T) {
	service.InitService()

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.GetUsersHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "[]\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateUserHandler(t *testing.T) {
	service.InitService()
	user := service.User{Name: "Test User", Email: "test@example.com"}
	jsonUser, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonUser))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.GetUsersHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createdUser service.User
	json.NewDecoder(rr.Body).Decode(&createdUser)

	if createdUser.Name != user.Name || createdUser.Email != user.Email || createdUser.ID == "" {
		t.Errorf("handler returned unexpected user: %+v", createdUser)
	}
}

func TestGetUserByIDHandler(t *testing.T) {
	service.InitService()
	// First, create a user
	user := service.User{Name: "Another User", Email: "another@example.com"}
	service.CreateUser(user) // This creates user with ID 1 (nextID was 1)

	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.GetUserByIDHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var fetchedUser service.User
	json.NewDecoder(rr.Body).Decode(&fetchedUser)

	if fetchedUser.ID != "1" || fetchedUser.Name != user.Name || fetchedUser.Email != user.Email {
		t.Errorf("handler returned unexpected user: %+v", fetchedUser)
	}
}

func TestGetUserByIDNotFoundHandler(t *testing.T) {
	service.InitService()
	req, err := http.NewRequest("GET", "/users/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.GetUserByIDHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}
