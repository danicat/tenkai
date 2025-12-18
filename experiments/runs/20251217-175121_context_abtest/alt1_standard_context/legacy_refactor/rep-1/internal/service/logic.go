package service

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserService defines the interface for user-related operations
type UserService interface {
	GetAllUsers() []User
	CreateUser(user User) (User, error)
	GetUserByID(id string) (User, error)
}

// inMemoryUserService implements UserService using an in-memory map
type inMemoryUserService struct {
	users map[string]User
	mutex sync.RWMutex
}

// NewInMemoryUserService creates a new instance of inMemoryUserService
func NewInMemoryUserService() UserService {
	return &inMemoryUserService{
		users: make(map[string]User),
	}
}

// GetAllUsers returns all users
func (s *inMemoryUserService) GetAllUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

// CreateUser creates a new user and assigns a UUID
func (s *inMemoryUserService) CreateUser(user User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user.ID = uuid.New().String()
	s.users[user.ID] = user
	return user, nil
}

// GetUserByID returns a user by ID
func (s *inMemoryUserService) GetUserByID(id string) (User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	u, ok := s.users[id]
	if !ok {
		return User{}, fmt.Errorf("user with ID %s not found", id)
	}
	return u, nil
}
