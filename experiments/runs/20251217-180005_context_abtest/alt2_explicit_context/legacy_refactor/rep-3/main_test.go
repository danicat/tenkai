package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"rep-3/internal/api"
	"rep-3/internal/service"
)

func TestUserAPI(t *testing.T) {
	// Initialize the service and handler for testing
	userService := service.NewInMemoryUserService()
	userHandler := api.NewUserHandler(userService)

	// Setup routes for the test server
	userHandler.SetupRoutes()

	t.Run("CreateUser", func(t *testing.T) {
		user := service.User{Name: "Alice", Email: "alice@example.com"}
		body, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var createdUser service.User
		json.NewDecoder(rr.Body).Decode(&createdUser)

		if createdUser.Name != user.Name || createdUser.Email != user.Email || createdUser.ID == "" {
			t.Errorf("handler returned unexpected body: got %+v want %+v with ID", createdUser, user)
		}
	})

	t.Run("ListUsers", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		rr := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var users []service.User
		json.NewDecoder(rr.Body).Decode(&users)

		if len(users) != 1 {
			t.Errorf("handler returned unexpected number of users: got %v want %v", len(users), 1)
		}

		if users[0].Name != "Alice" {
			t.Errorf("handler returned unexpected user: got %v want %v", users[0].Name, "Alice")
		}
	})

	t.Run("GetUserByID", func(t *testing.T) {
		// First, create a user to get an ID
		user := service.User{Name: "Bob", Email: "bob@example.com"}
		body, _ := json.Marshal(user)
		createReq, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		createReq.Header.Set("Content-Type", "application/json")
		createRr := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(createRr, createReq)

		var createdUser service.User
		json.NewDecoder(createRr.Body).Decode(&createdUser)
		userID := createdUser.ID

		// Now, get the user by ID
		getReq, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", userID), nil)
		getRr := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(getRr, getReq)

		if status := getRr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var fetchedUser service.User
		json.NewDecoder(getRr.Body).Decode(&fetchedUser)

		if fetchedUser.ID != userID || fetchedUser.Name != user.Name {
			t.Errorf("handler returned unexpected user: got %+v want %+v", fetchedUser, user)
		}
	})

	t.Run("GetUserNotFound", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/nonexistent", nil)
		rr := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})
}
