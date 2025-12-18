package service

import (
	"testing"
)

func TestNewUserService(t *testing.T) {
	svc := NewUserService()
	if svc == nil {
		t.Errorf("Expected NewUserService to return a non-nil service, got nil")
	}
	if svc.users == nil {
		t.Errorf("Expected users map to be initialized, got nil")
	}
	if svc.nextID != 1 {
		t.Errorf("Expected nextID to be 1, got %d", svc.nextID)
	}
}

func TestUserService_CreateUser(t *testing.T) {
	svc := NewUserService()
	user := User{Name: "Test User", Email: "test@example.com"}

	createdUser, err := svc.CreateUser(user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if createdUser.ID == "" {
		t.Errorf("Expected user ID to be generated, got empty string")
	}
	if createdUser.Name != user.Name || createdUser.Email != user.Email {
		t.Errorf("Created user mismatch. Expected %+v, got %+v", user, createdUser)
	}

	// Check if the user is actually stored
	storedUser, ok := svc.GetUserByID(createdUser.ID)
	if !ok {
		t.Errorf("Expected user to be stored, but not found")
	}
	if storedUser.Name != user.Name {
		t.Errorf("Stored user name mismatch. Expected %s, got %s", user.Name, storedUser.Name)
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	svc := NewUserService()
	user1 := User{Name: "User 1", Email: "user1@example.com"}
	user2 := User{Name: "User 2", Email: "user2@example.com"}

	svc.CreateUser(user1)
	svc.CreateUser(user2)

	users := svc.GetAllUsers()
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	svc := NewUserService()
	user := User{Name: "Find Me", Email: "findme@example.com"}

	createdUser, _ := svc.CreateUser(user)

	foundUser, ok := svc.GetUserByID(createdUser.ID)
	if !ok {
		t.Errorf("Expected to find user with ID %s", createdUser.ID)
	}
	if foundUser.Name != user.Name {
		t.Errorf("Found user name mismatch. Expected %s, got %s", user.Name, foundUser.Name)
	}

	_, ok = svc.GetUserByID("non-existent-id")
	if ok {
		t.Errorf("Expected not to find user with non-existent ID")
	}
}
