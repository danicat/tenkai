package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"experiment/internal/service"
)

func TestCreateAndGetUsers(t *testing.T) {
	// Clear existing users for a clean test environment (if necessary)
	// In a real application, you might use a mock or separate test database
	// For this example, we'll assume a fresh start or add a reset function to the service

	user1 := service.User{Name: "Alice", Email: "alice@example.com"}
	createdUser1 := service.CreateUser(user1)

	assert.NotEmpty(t, createdUser1.ID)
	assert.Equal(t, "Alice", createdUser1.Name)
	assert.Equal(t, "alice@example.com", createdUser1.Email)

	retrievedUser1, ok := service.GetUserByID(createdUser1.ID)
	assert.True(t, ok)
	assert.Equal(t, createdUser1, retrievedUser1)

	user2 := service.User{Name: "Bob", Email: "bob@example.com"}
	createdUser2 := service.CreateUser(user2)

	assert.NotEmpty(t, createdUser2.ID)

	allUsers := service.GetAllUsers()
	assert.Len(t, allUsers, 2)
	assert.Contains(t, allUsers, createdUser1)
	assert.Contains(t, allUsers, createdUser2)
}

func TestGetUserNotFound(t *testing.T) {
	_, ok := service.GetUserByID("nonexistent-id")
	assert.False(t, ok)
}
