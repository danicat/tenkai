package service_test

import (
	"testing"

	"experiment/internal/service"
)

func TestCreateUser(t *testing.T) {
	user := service.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	createdUser, err := service.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if createdUser.ID == "" {
		t.Error("Expected user ID to be generated, got empty string")
	}

	if createdUser.Name != user.Name {
		t.Errorf("Expected user name %s, got %s", user.Name, createdUser.Name)
	}
}

func TestGetAllUsers(t *testing.T) {
	service.CreateUser(service.User{Name: "User 1", Email: "user1@example.com"})
	service.CreateUser(service.User{Name: "User 2", Email: "user2@example.com"})

	users := service.GetAllUsers()
	if len(users) < 2 {
		t.Errorf("Expected at least 2 users, got %d", len(users))
	}
}

func TestGetUserByID(t *testing.T) {
	user := service.User{
		Name:  "Find Me",
		Email: "findme@example.com",
	}
	createdUser, _ := service.CreateUser(user)

	foundUser, ok := service.GetUserByID(createdUser.ID)
	if !ok {
		t.Fatalf("GetUserByID failed to find user with ID %s", createdUser.ID)
	}

	if foundUser.Name != user.Name {
		t.Errorf("Expected user name %s, got %s", user.Name, foundUser.Name)
	}

	_, ok = service.GetUserByID("nonexistent-id")
	if ok {
		t.Error("Expected GetUserByID to not find a nonexistent user")
	}
}
