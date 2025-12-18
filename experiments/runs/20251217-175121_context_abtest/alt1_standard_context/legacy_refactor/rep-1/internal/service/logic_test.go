package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryUserService_CreateUser(t *testing.T) {
	service := NewInMemoryUserService()

	user := User{Name: "Test User", Email: "test@example.com"}
	createdUser, err := service.CreateUser(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, createdUser.ID)
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)

	// Verify the user is in the storage
	retrievedUser, err := service.GetUserByID(createdUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdUser, retrievedUser)
}

func TestInMemoryUserService_GetAllUsers(t *testing.T) {
	service := NewInMemoryUserService()

	user1 := User{Name: "User 1", Email: "user1@example.com"}
	user2 := User{Name: "User 2", Email: "user2@example.com"}

	service.CreateUser(user1)
	service.CreateUser(user2)

	users := service.GetAllUsers()
	assert.Len(t, users, 2)

	// Check if both users are present (order doesn't matter for GetAll)
	found1 := false
	found2 := false
	for _, u := range users {
		if u.Name == user1.Name {
			found1 = true
		}
		if u.Name == user2.Name {
			found2 = true
		}
	}
	assert.True(t, found1)
	assert.True(t, found2)
}

func TestInMemoryUserService_GetUserByID(t *testing.T) {
	service := NewInMemoryUserService()

	user := User{Name: "Single User", Email: "single@example.com"}
	createdUser, _ := service.CreateUser(user)

	retrievedUser, err := service.GetUserByID(createdUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdUser, retrievedUser)

	// Test non-existent user
	_, err = service.GetUserByID("non-existent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
