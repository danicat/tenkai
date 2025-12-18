package service_test

import (
	"experiment/internal/service"
	"testing"
)

func TestUserService(t *testing.T) {
	svc := service.NewUserService()

	// Test CreateUser
	user1 := service.User{Name: "Alice", Email: "alice@example.com"}
	createdUser1, err := svc.CreateUser(user1)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if createdUser1.ID == "" {
		t.Error("Expected User ID to be generated, got empty")
	}

	// Test GetUserByID
	fetchedUser, ok := svc.GetUserByID(createdUser1.ID)
	if !ok {
		t.Fatalf("GetUserByID failed to find user with ID %s", createdUser1.ID)
	}
	if fetchedUser.Name != user1.Name || fetchedUser.Email != user1.Email {
		t.Errorf("Fetched user mismatch. Expected %+v, got %+v", user1, fetchedUser)
	}

	// Test GetAllUsers
	allUsers := svc.GetAllUsers()
	if len(allUsers) != 1 {
		t.Errorf("Expected 1 user, got %d", len(allUsers))
	}

	user2 := service.User{Name: "Bob", Email: "bob@example.com"}
	_, err = svc.CreateUser(user2)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	allUsers = svc.GetAllUsers()
	if len(allUsers) != 2 {
		t.Errorf("Expected 2 users, got %d", len(allUsers))
	}

	_, ok = svc.GetUserByID("nonexistent-id")
	if ok {
		t.Error("Expected GetUserByID to not find a nonexistent user")
	}
}
