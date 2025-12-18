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

// Service provides user-related business logic.
type Service struct {
	users  map[string]User
	mutex  sync.RWMutex
	nextID int
}

// NewService creates and returns a new Service instance.
func NewService() *Service {
	return &Service{
		users:  make(map[string]User),
		nextID: 1,
	}
}

// GetAllUsers returns all registered users.
func (s *Service) GetAllUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

// CreateUser creates a new user and adds it to the system.
func (s *Service) CreateUser(u User) User {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	u.ID = fmt.Sprintf("%d", s.nextID)
	s.nextID++
	s.users[u.ID] = u
	return u
}

// GetUserByID returns a user by their ID.
func (s *Service) GetUserByID(id string) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	u, ok := s.users[id]
	return u, ok
}
