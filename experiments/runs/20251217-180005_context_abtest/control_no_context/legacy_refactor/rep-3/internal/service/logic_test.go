package service

import (
	"reflect"
	"testing"
)

func TestNewUserService(t *testing.T) {
	service := NewUserService()
	if service == nil {
		t.Error("NewUserService returned nil")
	}
	if service.users == nil {
		t.Error("NewUserService did not initialize users map")
	}
	if service.nextID != 1 {
		t.Errorf("NewUserService nextID expected 1, got %d", service.nextID)
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	service := NewUserService()

	// Test with no users
	users := service.GetAllUsers()
	if len(users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(users))
	}

	// Test with some users
	user1 := User{Name: "Alice", Email: "alice@example.com"}
	user2 := User{Name: "Bob", Email: "bob@example.com"}

	service.CreateUser(user1)
	service.CreateUser(user2)

	users = service.GetAllUsers()
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	// Check if the users are present (order doesn't matter for map)
	foundAlice := false
	foundBob := false
	for _, u := range users {
		if u.Name == "Alice" {
			foundAlice = true
		}
		if u.Name == "Bob" {
			foundBob = true
		}
	}
	if !foundAlice || !foundBob {
		t.Errorf("Did not find expected users Alice or Bob")
	}
}

func TestUserService_CreateUser(t *testing.T) {
	service := NewUserService()

	user := User{Name: "Charlie", Email: "charlie@example.com"}
	createdUser := service.CreateUser(user)

	if createdUser.ID != "1" {
		t.Errorf("Expected ID '1', got '%s'", createdUser.ID)
	}
	if createdUser.Name != user.Name || createdUser.Email != user.Email {
		t.Errorf("Created user mismatch: expected %+v, got %+v", user, createdUser)
	}

	// Verify it's in the store
	storedUser, ok := service.users[createdUser.ID]
	if !ok || !reflect.DeepEqual(storedUser, createdUser) {
		t.Errorf("Created user not correctly stored or retrieved from map")
	}

	// Test incrementing ID
	user2 := User{Name: "David", Email: "david@example.com"}
	createdUser2 := service.CreateUser(user2)
	if createdUser2.ID != "2" {
		t.Errorf("Expected ID '2', got '%s'", createdUser2.ID)
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	service := NewUserService()

	user1 := User{Name: "Eve", Email: "eve@example.com"}
	createdUser1 := service.CreateUser(user1)

	// Test existing user
	retrievedUser, ok := service.GetUserByID(createdUser1.ID)
	if !ok {
		t.Errorf("Expected to find user with ID %s, but got not found", createdUser1.ID)
	}
	if !reflect.DeepEqual(retrievedUser, createdUser1) {
		t.Errorf("Retrieved user mismatch: expected %+v, got %+v", createdUser1, retrievedUser)
	}

	// Test non-existing user
	_, ok = service.GetUserByID("nonexistent-id")
	if ok {
		t.Errorf("Expected not to find user with non-existent ID, but found one")
	}
}
