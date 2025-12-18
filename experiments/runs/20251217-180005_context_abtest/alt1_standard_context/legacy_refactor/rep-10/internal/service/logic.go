package service

import (
	"sync"
	"github.com/google/uuid" // Will be used later for ID generation
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
	GetUserByID(id string) (User, bool)
	CreateUser(user User) (User, error)
}

// userService implements UserService
type userService struct {
	users  map[string]User
	mutex  sync.RWMutex
}

// NewUserService creates and returns a new UserService
func NewUserService() UserService {
	return &userService{
		users:  make(map[string]User),
		mutex:  sync.RWMutex{},
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

func (s *userService) GetUserByID(id string) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	u, ok := s.users[id]
	return u, ok
}

func (s *userService) CreateUser(user User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// For now, using simple integer ID. Will upgrade to uuid later.
	user.ID = uuid.New().String()
	s.users[user.ID] = user
	return user, nil
}
