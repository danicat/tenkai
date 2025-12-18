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

// UserService handles user-related business logic
type UserService struct {
	users map[string]User
	mutex sync.RWMutex
}

// NewUserService creates a new UserService
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]User),
	}
}

// CreateUser creates a new user and adds it to the system
func (s *UserService) CreateUser(u User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	u.ID = uuid.New().String()
	s.users[u.ID] = u
	return u, nil
}

// GetUsers returns all users
func (s *UserService) GetUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(id string) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	u, ok := s.users[id]
	return u, ok
}
