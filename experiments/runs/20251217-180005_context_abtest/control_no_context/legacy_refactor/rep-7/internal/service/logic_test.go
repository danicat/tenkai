package service

import (
	"testing"
)

func TestCreateAndGetUser(t *testing.T) {
	// Clear users map for a clean test environment
	users = make(map[string]User)

	user1 := User{Name: "Alice", Email: "alice@example.com"}
	createdUser1, err := CreateUser(user1)
	if err != nil {
		t.Fatalf("Error creating user 1: %v", err)
	}

	user2 := User{Name: "Bob", Email: "bob@example.com"}
	createdUser2, err := CreateUser(user2)
	if err != nil {
		t.Fatalf("Error creating user 2: %v", err)
	}

	if createdUser1.ID == createdUser2.ID {
		t.Errorf("Expected different IDs for created users, but got %s", createdUser1.ID)
	}

	retrievedUser1, ok := GetUserByID(createdUser1.ID)
	if !ok {
		t.Errorf("User with ID %s not found", createdUser1.ID)
	}

	if retrievedUser1.Name != user1.Name || retrievedUser1.Email != user1.Email {
		t.Errorf("Retrieved user 1 mismatch: expected %+v, got %+v", user1, retrievedUser1)
	}

	retrievedUser2, ok := GetUserByID(createdUser2.ID)
	if !ok {
		t.Errorf("User with ID %s not found", createdUser2.ID)
	}

	if retrievedUser2.Name != user2.Name || retrievedUser2.Email != user2.Email {
		t.Errorf("Retrieved user 2 mismatch: expected %+v, got %+v", user2, retrievedUser2)
	}

	_, ok = GetUserByID("nonexistent-id")
	if ok {
		t.Error("Expected no user for a nonexistent ID, but found one")
	}
}

func TestGetAllUsers(t *testing.T) {
	// Clear users map for a clean test environment
	users = make(map[string]User)

	list := GetAllUsers()
	if len(list) != 0 {
		t.Errorf("Expected empty list, got %d users", len(list))
	}

	user1 := User{Name: "Alice", Email: "alice@example.com"}
	CreateUser(user1)
	user2 := User{Name: "Bob", Email: "bob@example.com"}
	CreateUser(user2)

	list = GetAllUsers()
	if len(list) != 2 {
		t.Errorf("Expected 2 users, got %d users", len(list))
	}

	foundAlice := false
	foundBob := false
	for _, u := range list {
		if u.Name == "Alice" {
			foundAlice = true
		}
		if u.Name == "Bob" {
			foundBob = true
		}
	}

	if !foundAlice || !foundBob {
		t.Errorf("Expected to find Alice and Bob in the list, but didn't: %+v", list)
	}
}
