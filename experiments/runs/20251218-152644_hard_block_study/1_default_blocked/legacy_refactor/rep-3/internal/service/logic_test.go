package service

import (
	"testing"

	"github.com/google/uuid"
)

func TestCreateUser(t *testing.T) {
	service := NewUserService()
	user := User{Name: "Test User", Email: "test@example.com"}

	createdUser, err := service.CreateUser(user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if createdUser.Name != user.Name || createdUser.Email != user.Email {
		t.Errorf("Created user mismatch. Expected %+v, got %+v", user, createdUser)
	}

	if _, err := uuid.Parse(createdUser.ID); err != nil {
		t.Errorf("Expected a valid UUID for ID, got %s with error: %v", createdUser.ID, err)
	}
}

func TestGetUsers(t *testing.T) {
	service := NewUserService()
	user1 := User{Name: "User 1", Email: "user1@example.com"}
	user2 := User{Name: "User 2", Email: "user2@example.com"}

	service.CreateUser(user1)
	service.CreateUser(user2)

	users := service.GetUsers()
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

func TestGetUserByID(t *testing.T) {
	service := NewUserService()
	user := User{Name: "Test User", Email: "test@example.com"}
	createdUser, _ := service.CreateUser(user)

	// Test case: user found
	foundUser, ok := service.GetUserByID(createdUser.ID)
	if !ok {
		t.Errorf("Expected to find user with ID %s, but not found", createdUser.ID)
	}
	if foundUser.ID != createdUser.ID {
		t.Errorf("Expected user ID %s, got %s", createdUser.ID, foundUser.ID)
	}

	// Test case: user not found
	_, ok = service.GetUserByID("non-existent-id")
	if ok {
		t.Errorf("Expected not to find user with non-existent ID, but found")
	}
}
