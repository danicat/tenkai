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

type UserService struct {
	users  map[string]User
	mutex  sync.RWMutex
	nextID int // Only used if UUID generation fails or for fallback
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
	list := make([]User, 0, len(s.users))
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

func (s *UserService) CreateUser(user User) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Generate a new UUID for the user ID
	id, err := uuid.NewRandom()
	if err != nil {
		// Fallback to sequential ID if UUID generation fails
		// In a real application, you might want more robust error handling
		fmt.Printf("Warning: UUID generation failed, falling back to sequential ID: %v\n", err)
		user.ID = fmt.Sprintf("%d", s.nextID)
		s.nextID++
	} else {
		user.ID = id.String()
	}

	s.users[user.ID] = user
	return user, nil
}

func (s *UserService) GetUserByID(id string) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	u, ok := s.users[id]
	return u, ok
}
