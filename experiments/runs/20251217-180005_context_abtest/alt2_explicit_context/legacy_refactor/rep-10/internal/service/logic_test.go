package service_test

import (
	"experiment/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
	s := service.NewUserService()
	user, err := s.CreateUser("John Doe", "john@example.com")

	assert.Nil(t, err)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)

	users := s.GetUsers()
	assert.Len(t, users, 1)
}

func TestUserService_GetUsers(t *testing.T) {
	s := service.NewUserService()
	s.CreateUser("John Doe", "john@example.com")
	s.CreateUser("Jane Doe", "jane@example.com")

	users := s.GetUsers()
	assert.Len(t, users, 2)
}

func TestUserService_GetUser(t *testing.T) {
	s := service.NewUserService()
	createdUser, _ := s.CreateUser("John Doe", "john@example.com")

	user, err := s.GetUser(createdUser.ID)
	assert.Nil(t, err)
	assert.Equal(t, createdUser.ID, user.ID)
	assert.Equal(t, createdUser.Name, user.Name)

	_, err = s.GetUser("non-existent-id")
	assert.Equal(t, service.ErrUserNotFound, err)
}
