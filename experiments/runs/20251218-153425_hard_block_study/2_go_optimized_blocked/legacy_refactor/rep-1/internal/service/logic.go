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

// UserService provides operations for managing users.
type UserService struct {
	users map[string]User
	mutex sync.RWMutex
}

// NewUserService creates and returns a new UserService.
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]User),
	}
}

// CreateUser adds a new user to the system.
func (s *UserService) CreateUser(user User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user.ID = uuid.New().String()
	s.users[user.ID] = user
	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (s *UserService) GetUserByID(id string) (User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return User{}, fmt.Errorf("user with ID %s not found", id)
	}
	return user, nil
}

// GetAllUsers retrieves all users in the system.
func (s *UserService) GetAllUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}
