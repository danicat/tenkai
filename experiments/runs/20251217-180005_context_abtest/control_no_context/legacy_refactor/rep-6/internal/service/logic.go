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

// UserService defines the interface for user-related business logic
type UserService interface {
	CreateUser(user User) (User, error)
	GetUsers() []User
	GetUserByID(id string) (User, error)
}

type service struct {
	users map[string]User
	mutex sync.RWMutex
}

// NewService creates and returns a new UserService implementation
func NewService() UserService {
	return &service{
		users: make(map[string]User),
	}
}

func (s *service) CreateUser(user User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user.ID = uuid.New().String()
	s.users[user.ID] = user
	return user, nil
}

func (s *service) GetUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

func (s *service) GetUserByID(id string) (User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	u, ok := s.users[id]
	if !ok {
		return User{}, fmt.Errorf("user with ID %s not found", id)
	}
	return u, nil
}
