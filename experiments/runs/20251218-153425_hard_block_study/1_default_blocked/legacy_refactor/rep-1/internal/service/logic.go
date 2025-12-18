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
	users  map[string]User
	mutex  sync.RWMutex
}

// NewService creates and returns a new Service instance
func NewService() *Service {
	return &Service{
		users: make(map[string]User),
	}
}

// GetAllUsers returns a list of all users
func (s *Service) GetAllUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

// CreateUser creates a new user and adds it to the system
func (s *Service) CreateUser(u User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	newID, err := uuid.NewRandom()
	if err != nil {
		return User{}, fmt.Errorf("failed to generate UUID: %w", err)
	}
	u.ID = newID.String()
	s.users[u.ID] = u
	return u, nil
}

// GetUserByID retrieves a user by their ID
func (s *Service) GetUserByID(id string) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	u, ok := s.users[id]
	return u, ok
}

