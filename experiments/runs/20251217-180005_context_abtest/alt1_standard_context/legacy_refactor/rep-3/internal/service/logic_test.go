package service

import (
	"testing"
)

func TestGetAllUsers(t *testing.T) {
	// Clear users before test
	users = make(map[string]User)

	user1 := CreateUser("Alice", "alice@example.com")
	user2 := CreateUser("Bob", "bob@example.com")

	allUsers := GetAllUsers()

	if len(allUsers) != 2 {
		t.Errorf("Expected 2 users, got %d", len(allUsers))
	}

	found1 := false
	found2 := false
	for _, u := range allUsers {
		if u.ID == user1.ID && u.Name == user1.Name && u.Email == user1.Email {
			found1 = true
		}
		if u.ID == user2.ID && u.Name == user2.Name && u.Email == user2.Email {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Errorf("Did not find all created users in the list")
	}
}

func TestCreateUser(t *testing.T) {
	// Clear users before test
	users = make(map[string]User)

	name := "Charlie"
	email := "charlie@example.com"

	u := CreateUser(name, email)

	if u.Name != name {
		t.Errorf("Expected user name %s, got %s", name, u.Name)
	}
	if u.Email != email {
		t.Errorf("Expected user email %s, got %s", email, u.Email)
	}

	if _, ok := users[u.ID]; !ok {
		t.Errorf("User not found in map after creation")
	}
}

func TestGetUserByID(t *testing.T) {
	// Clear users before test
	users = make(map[string]User)

	user1 := CreateUser("David", "david@example.com")

	foundUser, ok := GetUserByID(user1.ID)

	if !ok {
		t.Errorf("Expected to find user with ID %s", user1.ID)
	}

	if foundUser.ID != user1.ID || foundUser.Name != user1.Name || foundUser.Email != user1.Email {
		t.Errorf("Found user mismatch. Expected %+v, got %+v", user1, foundUser)
	}

	_, ok = GetUserByID("nonexistent-id")
	if ok {
		t.Errorf("Expected not to find a user with a nonexistent ID")
	}
}
