package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"legacy_refactor/rep-2/internal/api"
	"legacy_refactor/rep-2/internal/service"
)

func TestUserFlow(t *testing.T) {
	// Initialize service and handlers for testing
	userService := service.NewUserService()
	handlers := api.NewUserHandlers(userService)

	// 1. Test Create User
	user1 := service.User{Name: "Alice", Email: "alice@example.com"}
	user1JSON, _ := json.Marshal(user1)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(user1JSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handlers.HandleUsers(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdUser1 service.User
	json.NewDecoder(rr.Body).Decode(&createdUser1)

	if createdUser1.Name != user1.Name || createdUser1.Email != user1.Email || createdUser1.ID == "" {
		t.Errorf("handler returned unexpected body: got %v want ID, Name and Email set", rr.Body.String())
	}

	// 2. Test Get All Users
	req, _ = http.NewRequest(http.MethodGet, "/users", nil)
	rr = httptest.NewRecorder()
	handlers.HandleUsers(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var users []service.User
	json.NewDecoder(rr.Body).Decode(&users)

	if len(users) != 1 {
		t.Errorf("handler returned wrong number of users: got %v want %v", len(users), 1)
	}

	// 3. Test Create another User
	user2 := service.User{Name: "Bob", Email: "bob@example.com"}
	user2JSON, _ := json.Marshal(user2)

	req, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(user2JSON))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handlers.HandleUsers(rr, req)

	var createdUser2 service.User
	json.NewDecoder(rr.Body).Decode(&createdUser2)

	// 4. Test Get User by ID
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", createdUser1.ID), nil)
	rr = httptest.NewRecorder()
	handlers.HandleUserByID(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for GetUserByID: got %v want %v", status, http.StatusOK)
	}

	var fetchedUser service.User
	json.NewDecoder(rr.Body).Decode(&fetchedUser)

	if fetchedUser.ID != createdUser1.ID || fetchedUser.Name != createdUser1.Name {
		t.Errorf("handler returned unexpected user by ID: got %v want %v", fetchedUser, createdUser1)
	}

	// 5. Test Get Non-Existent User by ID
	req, _ = http.NewRequest(http.MethodGet, "/users/nonexistent", nil)
	rr = httptest.NewRecorder()
	handlers.HandleUserByID(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for nonexistent user: got %v want %v", status, http.StatusNotFound)
	}

	// 6. Test invalid method for /users
	req, _ = http.NewRequest(http.MethodPut, "/users", nil)
	rr = httptest.NewRecorder()
	handlers.HandleUsers(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code for invalid method: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	// 7. Test invalid method for /users/id
	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/users/%s", createdUser1.ID), nil)
	rr = httptest.NewRecorder()
	handlers.HandleUserByID(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code for invalid method on ID: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	// 8. Test malformed JSON for Create User
	req, _ = http.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name": "Charlie", "email": "charlie@example.com", "extra": 123`))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handlers.HandleUsers(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for malformed JSON: got %v want %v", status, http.StatusBadRequest)
	}
}
