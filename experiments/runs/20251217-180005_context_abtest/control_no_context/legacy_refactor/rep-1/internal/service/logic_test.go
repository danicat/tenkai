package service

import (
	"testing"

	"github.com/google/uuid"
)

func TestCreateUser(t *testing.T) {
	// Clear users map before each test to ensure isolation
	users = make(map[uuid.UUID]User)

	name := "Test User"
	email := "test@example.com"

	user, err := CreateUser(name, email)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if user.Name != name {
		t.Errorf("Expected user name %s, got %s", name, user.Name)
	}
	if user.Email != email {
		t.Errorf("Expected user email %s, got %s", email, user.Email)
	}
	if user.ID == uuid.Nil {
		t.Errorf("Expected user ID to be non-nil, got %v", user.ID)
	}

	retrievedUser, ok := GetUserByID(user.ID)
	if !ok {
		t.Errorf("Expected user with ID %v to be found", user.ID)
	}
	if retrievedUser.Name != name {
		t.Errorf("Expected retrieved user name %s, got %s", name, retrievedUser.Name)
	}
}

func TestGetAllUsers(t *testing.T) {
	// Clear users map before each test
	users = make(map[uuid.UUID]User)

	user1, _ := CreateUser("User One", "one@example.com")
	user2, _ := CreateUser("User Two", "two@example.com")

	allUsers := GetAllUsers()

	if len(allUsers) != 2 {
		t.Errorf("Expected 2 users, got %d", len(allUsers))
	}

	found1 := false
	found2 := false
	for _, u := range allUsers {
		if u.ID == user1.ID {
			found1 = true
		}
		if u.ID == user2.ID {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Errorf("Expected to find both users in the list")
	}
}

func TestGetUserByID(t *testing.T) {
	// Clear users map before each test
	users = make(map[uuid.UUID]User)

	user, _ := CreateUser("Find Me", "findme@example.com")

	// Test case: user found
	retrievedUser, ok := GetUserByID(user.ID)
	if !ok {
		t.Errorf("Expected to find user with ID %v", user.ID)
	}
	if retrievedUser.Name != "Find Me" {
		t.Errorf("Expected user name 'Find Me', got %s", retrievedUser.Name)
	}

	// Test case: user not found
	nonExistentID := uuid.New()
	_, ok = GetUserByID(nonExistentID)
	if ok {
		t.Errorf("Expected not to find user with ID %v", nonExistentID)
	}
}
