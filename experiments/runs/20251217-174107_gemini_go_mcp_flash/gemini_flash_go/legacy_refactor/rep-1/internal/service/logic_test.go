package service

import (
	"testing"
)

func TestService_CreateUser(t *testing.T) {
	svc := NewService()

	u := User{Name: "Test User", Email: "test@example.com"}
	createdUser := svc.CreateUser(u)

	if createdUser.ID == "" {
		t.Error("Expected user ID to be set, got empty")
	}

	if createdUser.Name != u.Name || createdUser.Email != u.Email {
		t.Errorf("Expected user name and email to match, got %v", createdUser)
	}

	if createdUser.ID != "1" {
		t.Errorf("Expected first user ID to be \"1\", got %s", createdUser.ID)
	}

	// Original problematic line: u2 := User{Name: "Test User 2", Email: "test2@example.com"}
	// Fix: Directly pass the struct literal
	createdUser2 := svc.CreateUser(User{Name: "Test User 2", Email: "test2@example.com"})

	if createdUser2.ID != "2" {
		t.Errorf("Expected second user ID to be \"2\", got %s", createdUser2.ID)
	}
}

func TestService_GetAllUsers(t *testing.T) {
	svc := NewService()

	users := svc.GetAllUsers()
	if len(users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(users))
	}

	svc.CreateUser(User{Name: "User 1"})
	svc.CreateUser(User{Name: "User 2"})

	users = svc.GetAllUsers()
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

func TestService_GetUserByID(t *testing.T) {
	svc := NewService()

	u1 := svc.CreateUser(User{Name: "User 1", Email: "user1@example.com"})
	svc.CreateUser(User{Name: "User 2", Email: "user2@example.com"})

	// Test existing user
	foundUser, ok := svc.GetUserByID(u1.ID)
	if !ok {
		t.Errorf("Expected to find user with ID %s", u1.ID)
	}
	if foundUser.Name != u1.Name {
		t.Errorf("Expected user name %s, got %s", u1.Name, foundUser.Name)
	}

	// Test non-existing user
	_, ok = svc.GetUserByID("nonexistent")
	if ok {
		t.Error("Expected not to find a non-existent user")
	}
}