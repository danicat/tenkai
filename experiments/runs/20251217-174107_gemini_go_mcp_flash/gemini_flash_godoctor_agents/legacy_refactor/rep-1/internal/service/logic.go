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

// UserService handles user-related business logic
type UserService struct {
	users  map[string]User
	mutex  sync.RWMutex
	nextID int
}

// NewUserService creates and returns a new UserService
func NewUserService() *UserService {
	return &UserService{
		users:  make(map[string]User),
		nextID: 1,
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
func (s *UserService) CreateUser(u User) User {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	u.ID = fmt.Sprintf("%d", s.nextID)
	s.nextID++
	s.users[u.ID] = u
	return u
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(id string) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	u, ok := s.users[id]
	return u, ok
}
