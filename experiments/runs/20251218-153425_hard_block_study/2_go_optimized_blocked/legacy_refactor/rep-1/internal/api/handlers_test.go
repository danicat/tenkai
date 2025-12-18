package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"legacy_refactor/rep-1/internal/api"
	"legacy_refactor/rep-1/internal/service"
)

func TestUserHandlers(t *testing.T) {
	userService := service.NewUserService()
	handler := api.NewUserHandler(userService)

	var testUserID string

	t.Run("CreateUser", func(t *testing.T) {
		user := service.User{Name: "Test User", Email: "test@example.com"}
		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.HandleUsers(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var createdUser service.User
		json.NewDecoder(rr.Body).Decode(&createdUser)
		if createdUser.ID == "" {
			t.Error("expected user ID to be generated, got empty")
		}
		if createdUser.Name != user.Name {
			t.Errorf("expected user name %s, got %s", user.Name, createdUser.Name)
		}
		testUserID = createdUser.ID // Store the ID for subsequent tests
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rr := httptest.NewRecorder()

		handler.HandleUsers(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var users []service.User
		json.NewDecoder(rr.Body).Decode(&users)
		if len(users) != 1 {
			t.Errorf("expected 1 user, got %d", len(users))
		}
		if users[0].Name != "Test User" {
			t.Errorf("expected user name 'Test User', got %s", users[0].Name)
		}
	})

	t.Run("GetUserByID", func(t *testing.T) {
		if testUserID == "" {
			t.Fatal("testUserID not set, cannot run GetUserByID test")
		}

		req := httptest.NewRequest(http.MethodGet, "/users/"+testUserID, nil)
		rr := httptest.NewRecorder()

		handler.HandleUserByID(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var retrievedUser service.User
		json.NewDecoder(rr.Body).Decode(&retrievedUser)
		if retrievedUser.ID != testUserID {
			t.Errorf("expected user ID %s, got %s", testUserID, retrievedUser.ID)
		}
		if retrievedUser.Name != "Test User" {
			t.Errorf("expected user name 'Test User', got %s", retrievedUser.Name)
		}
	})
}
