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

// Service manages user-related business logic
type Service struct {
	users map[string]User
	mutex sync.RWMutex
}

// NewService creates a new instance of the UserService
func NewService() *Service {
	return &Service{
		users: make(map[string]User),
	}
}

// GetAllUsers retrieves all users
func (s *Service) GetAllUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

// CreateUser creates a new user
func (s *Service) CreateUser(user User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user.ID = uuid.New().String()
	s.users[user.ID] = user
	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(id string) (User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return User{}, fmt.Errorf("user with ID %s not found", id)
	}
	return u, nil
}
