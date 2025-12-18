package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"experiment/internal/service"
)

func TestGetUsersHandler(t *testing.T) {
	// Setup service and handlers
	svc := service.NewUserService()
	h := &Handlers{UserService: svc}

	// Create a few users
	svc.CreateUser(service.User{Name: "Alice", Email: "alice@example.com"})
	svc.CreateUser(service.User{Name: "Bob", Email: "bob@example.com"})

	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetUsersHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body.
	var users []service.User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("handler returned unexpected number of users: got %v want %v",
			len(users), 2)
	}

	// Further checks on user data can be added if necessary
}

func TestCreateUserHandler(t *testing.T) {
	// Setup service and handlers
	svc := service.NewUserService()
	h := &Handlers{UserService: svc}

	// Create a user to send in the request body
	newUser := service.User{Name: "Charlie", Email: "charlie@example.com"}
	body, _ := json.Marshal(newUser)

	// Create a request to pass to our handler.
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// We create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.CreateUserHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body.
	var createdUser service.User
	if err := json.NewDecoder(rr.Body).Decode(&createdUser); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	if createdUser.ID == "" {
		t.Errorf("Expected user ID to be generated, got empty string")
	}
	if createdUser.Name != newUser.Name || createdUser.Email != newUser.Email {
		t.Errorf("Created user mismatch. Expected %+v, got %+v", newUser, createdUser)
	}

	// Verify user is in the service
	_, ok := svc.GetUserByID(createdUser.ID)
	if !ok {
		t.Errorf("Expected user to be stored in service, but not found")
	}
}

func TestGetUserByIDHandler(t *testing.T) {
	// Setup service and handlers
	svc := service.NewUserService()
	h := &Handlers{UserService: svc}

	// Create a user to fetch
	userToFetch := service.User{Name: "David", Email: "david@example.com"}
	createdUser, _ := svc.CreateUser(userToFetch)

	// Test case 1: User found
	req, err := http.NewRequest("GET", "/users/"+createdUser.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetUserByIDHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for existing user: got %v want %v",
			status, http.StatusOK)
	}

	var fetchedUser service.User
	if err := json.NewDecoder(rr.Body).Decode(&fetchedUser); err != nil {
		t.Fatalf("Could not decode response for existing user: %v", err)
	}

	if fetchedUser.ID != createdUser.ID || fetchedUser.Name != createdUser.Name {
		t.Errorf("Fetched user mismatch. Expected %+v, got %+v", createdUser, fetchedUser)
	}

	// Test case 2: User not found
	req, err = http.NewRequest("GET", "/users/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent user: got %v want %v",
			status, http.StatusNotFound)
	}
}
