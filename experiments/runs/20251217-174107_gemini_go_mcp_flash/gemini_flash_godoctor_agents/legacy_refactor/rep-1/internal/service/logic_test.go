package service

import (
	"reflect"
	"testing"
)

func TestNewUserService(t *testing.T) {
	s := NewUserService()
	if s == nil {
		t.Errorf("NewUserService() returned nil")
	}
	if len(s.users) != 0 {
		t.Errorf("NewUserService() users map is not empty, got %d", len(s.users))
	}
	if s.nextID != 1 {
		t.Errorf("NewUserService() nextID is not 1, got %d", s.nextID)
	}
}

func TestUserService_CreateUser(t *testing.T) {
	s := NewUserService()
	user := User{Name: "Test User", Email: "test@example.com"}
	createdUser := s.CreateUser(user)

	if createdUser.ID != "1" {
		t.Errorf("CreateUser() ID not as expected, got %s", createdUser.ID)
	}
	if createdUser.Name != user.Name || createdUser.Email != user.Email {
		t.Errorf("CreateUser() user data mismatch, got %+v", createdUser)
	}
	if len(s.users) != 1 {
		t.Errorf("CreateUser() users map size not 1, got %d", len(s.users))
	}
	if s.users["1"].Name != user.Name {
		t.Errorf("CreateUser() user not found in map or data mismatch")
	}

	user2 := User{Name: "Another User", Email: "another@example.com"}
	createdUser2 := s.CreateUser(user2)
	if createdUser2.ID != "2" {
		t.Errorf("CreateUser() second user ID not as expected, got %s", createdUser2.ID)
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	s := NewUserService()
	user1 := User{Name: "User One", Email: "one@example.com"}
	user2 := User{Name: "User Two", Email: "two@example.com"}
	s.CreateUser(user1)
	s.CreateUser(user2)

	allUsers := s.GetAllUsers()
	if len(allUsers) != 2 {
		t.Errorf("GetAllUsers() returned %d users, expected 2", len(allUsers))
	}

	// Check if users are present (order might vary due to map iteration)
	foundUser1 := false
	foundUser2 := false
	for _, u := range allUsers {
		if u.Name == user1.Name && u.Email == user1.Email {
			foundUser1 = true
		}
		if u.Name == user2.Name && u.Email == user2.Email {
			foundUser2 = true
		}
	}
	if !foundUser1 || !foundUser2 {
		t.Errorf("GetAllUsers() did not return all created users")
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	s := NewUserService()
	user := User{Name: "Find Me", Email: "findme@example.com"}
	createdUser := s.CreateUser(user)

	foundUser, ok := s.GetUserByID(createdUser.ID)
	if !ok {
		t.Errorf("GetUserByID() failed to find user with ID %s", createdUser.ID)
	}
	if !reflect.DeepEqual(foundUser, createdUser) {
		t.Errorf("GetUserByID() returned wrong user. Got %+v, Expected %+v", foundUser, createdUser)
	}

	_, ok = s.GetUserByID("nonexistent")
	if ok {
		t.Errorf("GetUserByID() found a nonexistent user")
	}
}
