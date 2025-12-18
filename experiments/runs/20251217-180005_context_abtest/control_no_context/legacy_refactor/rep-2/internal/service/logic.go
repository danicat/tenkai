package service

import (
	"fmt"
	"sync"
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
	nextID int
}

// NewService creates and initializes a new Service
func NewService() *Service {
	return &Service{
		users:  make(map[string]User),
		nextID: 1,
	}
}

// CreateUser creates a new user
func (s *Service) CreateUser(u User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	u.ID = fmt.Sprintf("%d", s.nextID)
	s.nextID++
	s.users[u.ID] = u
	return u, nil
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(id string) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	u, ok := s.users[id]
	return u, ok
}

// ListUsers lists all users
func (s *Service) ListUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}
