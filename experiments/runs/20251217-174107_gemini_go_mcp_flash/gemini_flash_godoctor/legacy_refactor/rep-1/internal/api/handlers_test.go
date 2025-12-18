package api_test

import (
	"bytes"
	"encoding/json"
	"experiment/internal/api"
	"experiment/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandlers_GetAllUsers(t *testing.T) {
	svc := service.NewUserService()
	handlers := api.NewUserHandlers(svc)

	// Add some dummy users
	svc.CreateUser(service.User{Name: "Alice", Email: "alice@example.com"})
	svc.CreateUser(service.User{Name: "Bob", Email: "bob@example.com"})

	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handlers.GetAllUsers(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"id":"`
	if !bytes.Contains(rr.Body.Bytes(), []byte(expected)) {
		t.Errorf("handler returned unexpected body: got %v want containing %v",
			rr.Body.String(), expected)
	}

	var users []service.User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestUserHandlers_CreateUser(t *testing.T) {
	svc := service.NewUserService()
	handlers := api.NewUserHandlers(svc)

	user := service.User{Name: "Charlie", Email: "charlie@example.com"}
	jsonUser, _ := json.Marshal(user)

	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonUser))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handlers.CreateUser(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createdUser service.User
	if err := json.NewDecoder(rr.Body).Decode(&createdUser); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if createdUser.ID == "" {
		t.Error("expected user ID to be generated, got empty")
	}
	if createdUser.Name != user.Name {
		t.Errorf("expected user name %s, got %s", user.Name, createdUser.Name)
	}
}

func TestUserHandlers_GetUserByID(t *testing.T) {
	svc := service.NewUserService()
	handlers := api.NewUserHandlers(svc)

	user := service.User{Name: "David", Email: "david@example.com"}
	createdUser, _ := svc.CreateUser(user)

	req, err := http.NewRequest(http.MethodGet, "/users/"+createdUser.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handlers.GetUserByID(rr, req, createdUser.ID)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var foundUser service.User
	if err := json.NewDecoder(rr.Body).Decode(&foundUser); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if foundUser.ID != createdUser.ID {
		t.Errorf("expected user ID %s, got %s", createdUser.ID, foundUser.ID)
	}
}
