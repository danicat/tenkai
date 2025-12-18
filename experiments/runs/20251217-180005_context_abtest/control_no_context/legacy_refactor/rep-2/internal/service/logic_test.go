package service

import (
	"testing"
)

func TestCreateAndGetUser(t *testing.T) {
	s := NewService()
	user := User{Name: "Test User", Email: "test@example.com"}

	createdUser, err := s.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if createdUser.ID == "" {
		t.Errorf("Expected user ID to be generated, got empty")
	}
	if createdUser.Name != user.Name {
		t.Errorf("Expected user name %s, got %s", user.Name, createdUser.Name)
	}

	retrievedUser, ok := s.GetUser(createdUser.ID)
	if !ok {
		t.Errorf("GetUser failed to retrieve user with ID %s", createdUser.ID)
	}
	if retrievedUser.Name != user.Name {
		t.Errorf("Expected retrieved user name %s, got %s", user.Name, retrievedUser.Name)
	}
}

func TestListUsers(t *testing.T) {
	s := NewService()
	user1 := User{Name: "User 1", Email: "user1@example.com"}
	user2 := User{Name: "User 2", Email: "user2@example.com"}

	s.CreateUser(user1)
	s.CreateUser(user2)

	users := s.ListUsers()
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}
