package service_test

import (
	"testing"

	"legacy_refactor/internal/service"
)

func TestUserService_CreateUser(t *testing.T) {
	svc := service.NewUserService()

	user := service.User{Name: "John Doe", Email: "john@example.com"}
	createdUser, err := svc.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if createdUser.ID == "" {
		t.Error("Expected user ID to be generated, got empty string")
	}
	if createdUser.Name != user.Name {
		t.Errorf("Expected user name %s, got %s", user.Name, createdUser.Name)
	}

	// Try to retrieve the user
	retrievedUser, ok := svc.GetUserByID(createdUser.ID)
	if !ok {
		t.Error("Expected to retrieve created user, but got not found")
	}
	if retrievedUser.ID != createdUser.ID {
		t.Errorf("Expected retrieved user ID %s, got %s", createdUser.ID, retrievedUser.ID)
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	svc := service.NewUserService()

	user1 := service.User{Name: "John Doe", Email: "john@example.com"}
	user2 := service.User{Name: "Jane Doe", Email: "jane@example.com"}

	_, _ = svc.CreateUser(user1)
	_, _ = svc.CreateUser(user2)

	users := svc.GetAllUsers()
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	foundJohn := false
	foundJane := false
	for _, u := range users {
		if u.Name == "John Doe" {
			foundJohn = true
		}
		if u.Name == "Jane Doe" {
			foundJane = true
		}
	}

	if !foundJohn || !foundJane {
		t.Error("Expected to find both John and Jane Doe in the list")
	}
}