package service

import (
	"testing"

	"rep-1/internal/models"
)

func TestCreateAndGetUsers(t *testing.T) {
	// Clear existing users for a clean test environment
	users = make(map[string]models.User)

	// Test CreateUser
	user1 := models.User{Name: "Alice", Email: "alice@example.com"}
	createdUser1, err := CreateUser(user1)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if createdUser1.ID == "" {
		t.Error("Expected user ID to be generated, got empty")
	}
	if createdUser1.Name != user1.Name || createdUser1.Email != user1.Email {
		t.Errorf("Created user mismatch. Expected %+v, got %+v", user1, createdUser1)
	}

	user2 := models.User{Name: "Bob", Email: "bob@example.com"}
	_ , err = CreateUser(user2)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Test GetAllUsers
	allUsers := GetAllUsers()
	if len(allUsers) != 2 {
		t.Errorf("Expected 2 users, got %d", len(allUsers))
	}

	// Test GetUserByID
	fetchedUser1, ok := GetUserByID(createdUser1.ID)
	if !ok {
		t.Errorf("Expected to find user with ID %s, but not found", createdUser1.ID)
	}
	if fetchedUser1.Name != user1.Name {
		t.Errorf("Fetched user by ID mismatch. Expected name %s, got %s", user1.Name, fetchedUser1.Name)
	}

	_, ok = GetUserByID("nonexistent-id")
	if ok {
		t.Error("Expected not to find a nonexistent user, but found one")
	}
}
