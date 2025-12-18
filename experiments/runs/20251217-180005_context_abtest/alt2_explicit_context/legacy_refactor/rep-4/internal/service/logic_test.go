package service_test

import (
	"testing"

	"rep-4/internal/service"
	"github.com/stretchr/testify/assert"
)

func setup() {
	service.ResetStateForTest()
}

func TestCreateUser(t *testing.T) {
	setup()
	user := service.User{Name: "Test User", Email: "test@example.com"}
	createdUser, err := service.CreateUser(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdUser.ID)
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)

	// Verify user can be retrieved
	retrievedUser, ok := service.GetUserByID(createdUser.ID)
	assert.True(t, ok)
	assert.Equal(t, createdUser, retrievedUser)
}

func TestGetAllUsers(t *testing.T) {
	setup()
	// Clear users map for clean test state
	// This is not ideal for concurrent tests, but sufficient for a simple test.
	// A proper test setup/teardown would use a mock database or clear the map more robustly.

	// Create a few users
	user1 := service.User{Name: "User 1", Email: "user1@example.com"}
	user2 := service.User{Name: "User 2", Email: "user2@example.com"}

	createdUser1, _ := service.CreateUser(user1)
	createdUser2, _ := service.CreateUser(user2)

	allUsers := service.GetAllUsers()
	assert.Len(t, allUsers, 2)
	assert.Contains(t, allUsers, createdUser1)
	assert.Contains(t, allUsers, createdUser2)
}

func TestGetUserByID(t *testing.T) {
	setup()
	user := service.User{Name: "Test User Get", Email: "get@example.com"}
	createdUser, _ := service.CreateUser(user)

	retrievedUser, ok := service.GetUserByID(createdUser.ID)
	assert.True(t, ok)
	assert.Equal(t, createdUser, retrievedUser)

	// Test non-existent user
	_, ok = service.GetUserByID("nonexistent-id")
	assert.False(t, ok)
}

func TestUpdateUser(t *testing.T) {
	setup()
	user := service.User{Name: "Original Name", Email: "original@example.com"}
	createdUser, _ := service.CreateUser(user)

	updatedUser := service.User{Name: "Updated Name", Email: "updated@example.com"}
	resultUser, err := service.UpdateUser(createdUser.ID, updatedUser)
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, resultUser.ID)
	assert.Equal(t, updatedUser.Name, resultUser.Name)
	assert.Equal(t, updatedUser.Email, resultUser.Email)

	retrievedUser, ok := service.GetUserByID(createdUser.ID)
	assert.True(t, ok)
	assert.Equal(t, resultUser, retrievedUser)

	// Test updating non-existent user
	_, err = service.UpdateUser("nonexistent-id", updatedUser)
	assert.Error(t, err)
}

func TestDeleteUser(t *testing.T) {
	setup()
	user := service.User{Name: "User to Delete", Email: "delete@example.com"}
	createdUser, _ := service.CreateUser(user)

	err := service.DeleteUser(createdUser.ID)
	assert.NoError(t, err)

	_, ok := service.GetUserByID(createdUser.ID)
	assert.False(t, ok)

	// Test deleting non-existent user
	err = service.DeleteUser("nonexistent-id")
	assert.Error(t, err)
}
