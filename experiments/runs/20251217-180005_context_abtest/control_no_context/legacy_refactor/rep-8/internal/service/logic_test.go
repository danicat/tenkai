package service_test

import (
	"rep-8/internal/service"
	"testing"

	"github.com/google/uuid"
)

func TestUserService_CreateUser(t *testing.T) {
	svc := service.NewUserService()

	user := service.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	createdUser := svc.CreateUser(user)

	if createdUser.ID == "" {
		t.Errorf("Expected user ID to be generated, got empty string")
	}

	// Verify if the generated ID is a valid UUID
	if _, err := uuid.Parse(createdUser.ID); err != nil {
		t.Errorf("Generated ID is not a valid UUID: %v", err)
	}

	if createdUser.Name != user.Name {
		t.Errorf("Expected user name %s, got %s", user.Name, createdUser.Name)
	}

	if createdUser.Email != user.Email {
		t.Errorf("Expected user email %s, got %s", user.Email, createdUser.Email)
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	svc := service.NewUserService()

	user1 := service.User{Name: "User 1", Email: "user1@example.com"}
	user2 := service.User{Name: "User 2", Email: "user2@example.com"}

	svc.CreateUser(user1)
	svc.CreateUser(user2)

	allUsers := svc.GetAllUsers()

	if len(allUsers) != 2 {
		t.Errorf("Expected 2 users, got %d", len(allUsers))
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	svc := service.NewUserService()

	user := service.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	createdUser := svc.CreateUser(user)

	retrievedUser, ok := svc.GetUserByID(createdUser.ID)

	if !ok {
		t.Errorf("Expected to find user with ID %s", createdUser.ID)
	}

	if retrievedUser.ID != createdUser.ID {
		t.Errorf("Expected user ID %s, got %s", createdUser.ID, retrievedUser.ID)
	}

	_, ok = svc.GetUserByID("non-existent-id")
	if ok {
		t.Errorf("Expected not to find user with non-existent ID")
	}
}
