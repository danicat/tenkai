package service

import (
	"sync"

	"github.com/google/uuid" // Will be updated to latest version later
)

// User represents a user in the system
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserService struct {
	users map[string]User
	mutex sync.RWMutex
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]User),
	}
}

func (s *UserService) GetAllUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

func (s *UserService) CreateUser(user User) User {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	user.ID = uuid.New().String()
	s.users[user.ID] = user
	return user
}

func (s *UserService) GetUserByID(id string) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	user, ok := s.users[id]
	return user, ok
}
