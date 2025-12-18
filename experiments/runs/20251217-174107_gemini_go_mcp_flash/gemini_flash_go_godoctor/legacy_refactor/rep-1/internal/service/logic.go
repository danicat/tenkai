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

// Service provides user-related business logic
type Service struct {
	users map[string]User
	mutex sync.RWMutex
}

// NewService creates and returns a new Service instance
func NewService() *Service {
	return &Service{
		users: make(map[string]User),
	}
}

// CreateUser adds a new user to the system
func (s *Service) CreateUser(user User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user.ID = uuid.New().String()
	s.users[user.ID] = user
	return user, nil
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(id string) (User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return User{}, fmt.Errorf("user with ID %s not found", id)
	}
	return user, nil
}

// ListUsers retrieves all users
func (s *Service) ListUsers() ([]User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list, nil
}
