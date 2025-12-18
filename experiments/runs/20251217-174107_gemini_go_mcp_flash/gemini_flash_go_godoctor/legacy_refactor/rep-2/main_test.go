package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"legacy_refactor/rep-2/internal/api"
	"legacy_refactor/rep-2/internal/service"
)

func TestUserAPI(t *testing.T) {
	svc := service.NewService()
	handler := api.NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListUsersHandler(w, r)
		case http.MethodPost:
			handler.CreateUserHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetUserHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := ts.Client()

	// Test CreateUser (POST /users)
	newUser := service.User{Name: "Alice", Email: "alice@example.com"}
	userJSON, _ := json.Marshal(newUser)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(res.Body)
		t.Fatalf("Expected status %d, got %d: %s", http.StatusCreated, res.StatusCode, string(b))
	}

	var createdUser service.User
	json.NewDecoder(res.Body).Decode(&createdUser)
	if createdUser.ID == "" || createdUser.Name != newUser.Name || createdUser.Email != newUser.Email {
		t.Errorf("Unexpected created user: %+v", createdUser)
	}

	// Test ListUsers (GET /users)
	req, _ = http.NewRequest(http.MethodGet, ts.URL+"/users", nil)
	res, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		t.Fatalf("Expected status %d, got %d: %s", http.StatusOK, res.StatusCode, string(b))
	}

	var users []service.User
	json.NewDecoder(res.Body).Decode(&users)
	if len(users) != 1 || users[0].ID != createdUser.ID {
		t.Errorf("Unexpected user list: %+v", users)
	}

	// Test GetUser (GET /users/{id})
	req, _ = http.NewRequest(http.MethodGet, ts.URL+"/users/"+createdUser.ID, nil)
	res, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		t.Fatalf("Expected status %d, got %d: %s", http.StatusOK, res.StatusCode, string(b))
	}

	var fetchedUser service.User
	json.NewDecoder(res.Body).Decode(&fetchedUser)
	if fetchedUser.ID != createdUser.ID || fetchedUser.Name != createdUser.Name || fetchedUser.Email != createdUser.Email {
		t.Errorf("Unexpected fetched user: %+v", fetchedUser)
	}
}
