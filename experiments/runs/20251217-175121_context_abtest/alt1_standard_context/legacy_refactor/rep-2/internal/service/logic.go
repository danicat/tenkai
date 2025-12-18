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

// UserService provides methods for interacting with users
type UserService struct {
	users  map[string]User
	mutex  sync.RWMutex
}

// NewUserService creates and returns a new UserService
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]User),
	}
}

// GetAllUsers returns a list of all users
func (s *UserService) GetAllUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

// CreateUser creates a new user and adds it to the system
func (s *UserService) CreateUser(user User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user.ID = uuid.New().String()
	s.users[user.ID] = user
	return user, nil
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id string) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	u, ok := s.users[id]
	return u, ok
}
