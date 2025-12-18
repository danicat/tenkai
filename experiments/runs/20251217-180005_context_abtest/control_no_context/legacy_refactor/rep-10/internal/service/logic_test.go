package service_test

import (
	"testing"

	"experiment/internal/service"
)

func TestCreateUser(t *testing.T) {
	u := service.User{Name: "Test User", Email: "test@example.com"}
	createdUser, err := service.CreateUser(u)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if createdUser.ID == "" {
		t.Error("expected user ID to be generated, got empty string")
	}
	if createdUser.Name != u.Name {
		t.Errorf("expected user name %q, got %q", u.Name, createdUser.Name)
	}
	if createdUser.Email != u.Email {
		t.Errorf("expected user email %q, got %q", u.Email, createdUser.Email)
	}
}

func TestGetAllUsers(t *testing.T) {
	// Clear existing users for a clean test environment
	service.ClearUsers()

	u1 := service.User{Name: "User 1", Email: "user1@example.com"}
	u2 := service.User{Name: "User 2", Email: "user2@example.com"}

	service.CreateUser(u1)
	service.CreateUser(u2)

	allUsers := service.GetAllUsers()
	if len(allUsers) != 2 {
		t.Errorf("expected 2 users, got %d", len(allUsers))
	}
}

func TestGetUserByID(t *testing.T) {
	// Clear existing users for a clean test environment
	service.ClearUsers()

	u := service.User{Name: "Find Me", Email: "findme@example.com"}
	createdUser, _ := service.CreateUser(u)

	foundUser, ok := service.GetUserByID(createdUser.ID)
	if !ok {
		t.Error("expected to find user, but got not found")
	}
	if foundUser.ID != createdUser.ID {
		t.Errorf("expected user ID %q, got %q", createdUser.ID, foundUser.ID)
	}

	_, ok = service.GetUserByID("nonexistent-id")
	if ok {
		t.Error("expected not to find user, but found one")
	}
}
