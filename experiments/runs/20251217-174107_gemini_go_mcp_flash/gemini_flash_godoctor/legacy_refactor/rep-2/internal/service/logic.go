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

// UserService defines the interface for user-related business logic
type UserService interface {
	GetAllUsers() []User
	CreateUser(u User) (User, error)
	GetUserByID(id string) (User, bool)
}

type userService struct {
	users map[string]User
	mutex sync.RWMutex
}

// NewUserService creates a new instance of UserService
func NewUserService() UserService {
	return &userService{
		users: make(map[string]User),
	}
}

func (s *userService) GetAllUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

func (s *userService) CreateUser(u User) (User, error) {
	    s.mutex.Lock()
	    defer s.mutex.Unlock()
	// Generate a UUID for the user ID
	u.ID = uuid.New().String()
	s.users[u.ID] = u
	return u, nil
}

func (s *userService) GetUserByID(id string) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	u, ok := s.users[id]
	return u, ok
}
