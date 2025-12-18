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
	// Initialize the UserService
	userService := service.NewUserService()

	// Initialize the UserHandler with the UserService
	userHandler := api.NewUserHandler(userService)

	// Create a new ServeMux to register the handlers
	mux := http.NewServeMux()
	mux.Handle("/users", userHandler)
	mux.Handle("/users/", userHandler)

	// Test POST /users
	t.Run("CreateUser", func(t *testing.T) {
		newUser := service.User{Name: "Alice", Email: "alice@example.com"}
		body, _ := json.Marshal(newUser)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var createdUser service.User
		json.NewDecoder(rr.Body).Decode(&createdUser)
		if createdUser.Name != newUser.Name || createdUser.Email != newUser.Email || createdUser.ID == "" {
			t.Errorf("handler returned unexpected body: got %v", createdUser)
		}
	})

	// Test GET /users
	t.Run("GetAllUsers", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var users []service.User
		json.NewDecoder(rr.Body).Decode(&users)
		if len(users) != 1 {
			t.Errorf("handler returned unexpected number of users: got %v want %v", len(users), 1)
		}
	})

	// Get a user ID to test GET /users/{id}
	var testUserID string
	// First, create a user to ensure there's at least one
	newUser := service.User{Name: "Bob", Email: "bob@example.com"}
	body, _ := json.Marshal(newUser)
	reqCreate, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	rrCreate := httptest.NewRecorder()
	mux.ServeHTTP(rrCreate, reqCreate)
	var createdUser service.User
	json.NewDecoder(rrCreate.Body).Decode(&createdUser)
	testUserID = createdUser.ID


	// Test GET /users/{id}
	t.Run("GetUserByID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/"+testUserID, nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var user service.User
		json.NewDecoder(rr.Body).Decode(&user)
		if user.ID != testUserID || user.Name != newUser.Name || user.Email != newUser.Email {
			t.Errorf("handler returned unexpected user: got %v", user)
		}
	})

	// Test GET /users/{id} with non-existent ID
	t.Run("GetUserByID_NotFound", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/nonexistent-id", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code for not found: got %v want %v", status, http.StatusNotFound)
		}
	})
}
