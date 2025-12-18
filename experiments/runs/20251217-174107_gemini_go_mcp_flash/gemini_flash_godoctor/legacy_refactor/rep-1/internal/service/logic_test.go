package service_test

import (
	"experiment/internal/service"
	"testing"

	"github.com/google/uuid"
)

func TestUserService_CreateUser(t *testing.T) {
	s := service.NewUserService()

	user := service.User{Name: "Test User", Email: "test@example.com"}
	createdUser, err := s.CreateUser(user)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if createdUser.ID == "" {
		t.Error("expected user ID to be generated, got empty")
	}

	if _, err := uuid.Parse(createdUser.ID); err != nil {
		t.Errorf("expected a valid UUID, got %v", createdUser.ID)
	}

	if createdUser.Name != user.Name {
		t.Errorf("expected user name %s, got %s", user.Name, createdUser.Name)
	}
	
	if createdUser.Email != user.Email {
		t.Errorf("expected user email %s, got %s", user.Email, createdUser.Email)
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	s := service.NewUserService()

	user1 := service.User{Name: "User 1", Email: "user1@example.com"}
	user2 := service.User{Name: "User 2", Email: "user2@example.com"}

	s.CreateUser(user1)
	s.CreateUser(user2)

	allUsers := s.GetAllUsers()

	if len(allUsers) != 2 {
		t.Errorf("expected 2 users, got %d", len(allUsers))
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	s := service.NewUserService()

	user := service.User{Name: "Test User", Email: "test@example.com"}
	createdUser, _ := s.CreateUser(user)

	foundUser, ok := s.GetUserByID(createdUser.ID)

	if !ok {
		t.Error("expected to find user, but got not found")
	}

	if foundUser.ID != createdUser.ID {
		t.Errorf("expected user ID %s, got %s", createdUser.ID, foundUser.ID)
	}

	_, ok = s.GetUserByID("nonexistent-id")
	if ok {
		t.Error("expected not to find user, but found one")
	}
}
