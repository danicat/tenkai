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

type UserService struct {
	users  map[string]User
	mutex  sync.RWMutex
	nextID int // Used for sequential IDs when not using UUIDs
}

func NewUserService() *UserService {
	return &UserService{
		users:  make(map[string]User),
		nextID: 1,
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

func (s *UserService) CreateUser(user User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Use UUID for new user IDs
	user.ID = uuid.New().String()
	s.users[user.ID] = user
	return user, nil
}

func (s *UserService) GetUserByID(id string) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	user, ok := s.users[id]
	return user, ok
}
