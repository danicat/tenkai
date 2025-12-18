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

// UserService provides business logic for users
type UserService struct {
	users  map[string]User
	mutex  sync.RWMutex
}

// NewUserService creates a new UserService
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]User),
	}
}

// GetAllUsers returns all users
func (s *UserService) GetAllUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

// CreateUser creates a new user and returns it
func (s *UserService) CreateUser(name, email string) (*User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}

	u := User{
		ID:    id.String(),
		Name:  name,
		Email: email,
	}
	s.users[u.ID] = u
	return &u, nil
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(id string) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	u, ok := s.users[id]
	if !ok {
		return nil, fmt.Errorf("user with ID %s not found", id)
	}
	return &u, nil
}
