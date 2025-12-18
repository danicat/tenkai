package service

import (
	"sync"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Service handles business logic related to users
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

// CreateUser creates a new user
func (s *Service) CreateUser(u User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	u.ID = uuid.New().String()
	s.users[u.ID] = u
	return u, nil
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(id string) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	user, ok := s.users[id]
	return user, ok
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
