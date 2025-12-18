package service

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

// User represents a user in the system
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserServicer interface {
	CreateUser(name, email string) (User, error)
	GetUsers() []User
	GetUser(id string) (User, error)
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

func (s *UserService) CreateUser(name, email string) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	u := User{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
	}
	s.users[u.ID] = u
	return u, nil
}

func (s *UserService) GetUsers() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var list []User
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

func (s *UserService) GetUser(id string) (User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	u, ok := s.users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}
