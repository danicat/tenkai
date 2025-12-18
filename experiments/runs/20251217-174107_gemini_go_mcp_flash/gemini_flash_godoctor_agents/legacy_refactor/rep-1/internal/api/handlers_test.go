package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"experiment/internal/service"
)

func TestUserHandler_HandleUsers_GET(t *testing.T) {
	userService := service.NewUserService()
	handler := NewUserHandler(userService)

	// Add some initial users
	userService.CreateUser(service.User{Name: "Alice", Email: "alice@example.com"})
	userService.CreateUser(service.User{Name: "Bob", Email: "bob@example.com"})

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HandleUsers(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var users []service.User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
	// Further checks for user data can be added
}

func TestUserHandler_HandleUsers_POST(t *testing.T) {
	userService := service.NewUserService()
	handler := NewUserHandler(userService)

	userJSON := `{"Name": "Charlie", "Email": "charlie@example.com"}`
	req, err := http.NewRequest("POST", "/users", strings.NewReader(userJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.HandleUsers(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createdUser service.User
	if err := json.NewDecoder(rr.Body).Decode(&createdUser); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if createdUser.ID == "" {
		t.Errorf("expected user ID to be set, got empty")
	}
	if createdUser.Name != "Charlie" || createdUser.Email != "charlie@example.com" {
		t.Errorf("created user data mismatch: got %+v", createdUser)
	}

	// Verify it's in the service
	_, ok := userService.GetUserByID(createdUser.ID)
	if !ok {
		t.Errorf("user not found in service after creation")
	}
}

func TestUserHandler_HandleUserByID_GET(t *testing.T) {
	userService := service.NewUserService()
	handler := NewUserHandler(userService)

	user := userService.CreateUser(service.User{Name: "David", Email: "david@example.com"})

	req, err := http.NewRequest("GET", "/users/"+user.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HandleUserByID(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var foundUser service.User
	if err := json.NewDecoder(rr.Body).Decode(&foundUser); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if foundUser.ID != user.ID || foundUser.Name != user.Name {
		t.Errorf("found user data mismatch: got %+v, expected %+v", foundUser, user)
	}
}

func TestUserHandler_HandleUserByID_GET_NotFound(t *testing.T) {
	userService := service.NewUserService()
	handler := NewUserHandler(userService)

	req, err := http.NewRequest("GET", "/users/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HandleUserByID(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for not found: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestUserHandler_HandleUsers_MethodNotAllowed(t *testing.T) {
	userService := service.NewUserService()
	handler := NewUserHandler(userService)

	req, err := http.NewRequest("PUT", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HandleUsers(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code for method not allowed: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestUserHandler_HandleUserByID_MethodNotAllowed(t *testing.T) {
	userService := service.NewUserService()
	handler := NewUserHandler(userService)

	req, err := http.NewRequest("POST", "/users/1", nil) // POST to /users/id is not allowed
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HandleUserByID(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code for method not allowed: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}
