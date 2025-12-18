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

// UserService defines the interface for user-related operations
type UserService interface {
	GetAllUsers() []User
	CreateUser(u User) (User, error)
	GetUserByID(id string) (User, bool)
}

// inMemoryUserService is an in-memory implementation of UserService
type inMemoryUserService struct {
	users map[string]User
	mutex sync.RWMutex
}

// NewInMemoryUserService creates a new in-memory user service
func NewInMemoryUserService() UserService {
	return &inMemoryUserService{
		users: make(map[string]User),
	}
}

// GetAllUsers retrieves all users
func (s *inMemoryUserService) GetAllUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

// CreateUser creates a new user and assigns a unique ID
func (s *inMemoryUserService) CreateUser(u User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	u.ID = uuid.New().String()
	s.users[u.ID] = u
	return u, nil
}

// GetUserByID retrieves a user by ID
func (s *inMemoryUserService) GetUserByID(id string) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	u, ok := s.users[id]
	return u, ok
}
