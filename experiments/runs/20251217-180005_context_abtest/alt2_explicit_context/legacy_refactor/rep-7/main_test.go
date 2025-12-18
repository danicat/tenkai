package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"rep-7/internal/api"
	"rep-7/internal/service"
)

func TestGetUsers(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.GetUsersHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK)
	}

	expected := "null\n"
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rec.Body.String(), expected)
	}
}

func TestCreateUser(t *testing.T) {
	user := service.User{Name: "Test User", Email: "test@example.com"}
	jsonBody, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			api.CreateUserHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status,
			http.StatusCreated)
	}

	var createdUser service.User
	json.NewDecoder(rec.Body).Decode(&createdUser)
	if createdUser.Name != user.Name || createdUser.Email != user.Email || createdUser.ID == "" {
		t.Errorf("handler returned unexpected user: got %+v want %+v with an ID",
			createdUser,
			user)
	}

	// Test GetUser by ID
	req, err = http.NewRequest("GET", "/users/"+createdUser.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	rec = httptest.NewRecorder()
	handler = http.HandlerFunc(api.GetUserHandler)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for GetUser: got %v want %v",
			status,
			http.StatusOK)
	}

	var fetchedUser service.User
	json.NewDecoder(rec.Body).Decode(&fetchedUser)
	if fetchedUser.ID != createdUser.ID || fetchedUser.Name != createdUser.Name || fetchedUser.Email != createdUser.Email {
		t.Errorf("handler returned unexpected fetched user: got %+v want %+v",
			fetchedUser,
			createdUser)
	}
}
