package service_test

import (
	"legacy_refactor/rep-1/internal/service"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	svc := service.NewUserService()

	user := service.User{Name: "John Doe", Email: "john@example.com"}
	createdUser, err := svc.CreateUser(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if createdUser.ID == "" {
		t.Error("expected user ID to be generated, got empty")
	}
	if createdUser.Name != user.Name {
		t.Errorf("expected name %s, got %s", user.Name, createdUser.Name)
	}

	retrievedUser, err := svc.GetUserByID(createdUser.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrievedUser.ID != createdUser.ID {
		t.Errorf("expected user ID %s, got %s", createdUser.ID, retrievedUser.ID)
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	svc := service.NewUserService()

	_, err := svc.GetUserByID("nonexistent-id")
	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}

	user := service.User{Name: "Jane Doe", Email: "jane@example.com"}
	createdUser, _ := svc.CreateUser(user)

	retrievedUser, err := svc.GetUserByID(createdUser.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrievedUser.Name != user.Name {
		t.Errorf("expected user name %s, got %s", user.Name, retrievedUser.Name)
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	svc := service.NewUserService()

	users := svc.GetAllUsers()
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}

	user1 := service.User{Name: "Alice", Email: "alice@example.com"}
	user2 := service.User{Name: "Bob", Email: "bob@example.com"}

	svc.CreateUser(user1)
	svc.CreateUser(user2)

	allUsers := svc.GetAllUsers()
	if len(allUsers) != 2 {
		t.Errorf("expected 2 users, got %d", len(allUsers))
	}
}
