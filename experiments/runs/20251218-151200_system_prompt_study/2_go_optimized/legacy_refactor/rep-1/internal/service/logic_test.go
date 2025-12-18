package service_test

import (
	"testing"

	"rep-1/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestNewUserService(t *testing.T) {
	s := service.NewUserService()
	assert.NotNil(t, s)
}

func TestCreateUser(t *testing.T) {
	s := service.NewUserService()

	user := service.User{Name: "John Doe", Email: "john@example.com"}
	createdUser, err := s.CreateUser(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdUser.ID)
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)

	// Verify user is in the system
	retrievedUser, ok := s.GetUserByID(createdUser.ID)
	assert.True(t, ok)
	assert.Equal(t, createdUser, retrievedUser)
}

func TestGetAllUsers(t *testing.T) {
	s := service.NewUserService()

	user1 := service.User{Name: "John Doe", Email: "john@example.com"}
	user2 := service.User{Name: "Jane Smith", Email: "jane@example.com"}

	s.CreateUser(user1)
	s.CreateUser(user2)

	allUsers := s.GetAllUsers()
	assert.Len(t, allUsers, 2)

	// Check if users are present (order doesn't matter for GetAll)
	foundJohn := false
	foundJane := false
	for _, u := range allUsers {
		if u.Name == "John Doe" {
			foundJohn = true
		}
		if u.Name == "Jane Smith" {
			foundJane = true
		}
	}
	assert.True(t, foundJohn)
	assert.True(t, foundJane)
}

func TestGetUserByID(t *testing.T) {
	s := service.NewUserService()

	user := service.User{Name: "John Doe", Email: "john@example.com"}
	createdUser, _ := s.CreateUser(user)

	// Test existing user
	retrievedUser, ok := s.GetUserByID(createdUser.ID)
	assert.True(t, ok)
	assert.Equal(t, createdUser, retrievedUser)

	// Test non-existing user
	_, ok = s.GetUserByID("non-existent-id")
	assert.False(t, ok)
}
