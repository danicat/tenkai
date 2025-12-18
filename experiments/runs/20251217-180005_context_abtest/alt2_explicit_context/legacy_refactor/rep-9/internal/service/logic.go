package service

import (
	"fmt"
	"sync"

	"github.com/google/uuid" // Will be updated to latest version later
)

// User represents a user in the system
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserService defines the interface for user-related business logic
type UserService interface {
	GetAllUsers() []User
	CreateUser(user User) (User, error)
	GetUserByID(id string) (User, error)
}

type userService struct {
	users  map[string]User
	mutex  sync.RWMutex
	// nextID is removed as we will use UUIDs
}

// NewUserService creates and returns a new UserService
func NewUserService() UserService {
	return &userService{
		users: make(map[string]User),
	}
}

func (s *userService) GetAllUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

func (s *userService) CreateUser(user User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Generate a new UUID for the user ID
	user.ID = uuid.New().String()
	s.users[user.ID] = user
	return user, nil
}

func (s *userService) GetUserByID(id string) (User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	u, ok := s.users[id]
	if !ok {
		return User{}, fmt.Errorf("user with ID %s not found", id)
	}
	return u, nil
}
