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

var (
	users = make(map[string]User)
	mutex sync.RWMutex
)

// GetAllUsers returns all users
func GetAllUsers() []User {
	mutex.RLock()
	defer mutex.RUnlock()
	var list []User
	for _, u := range users {
		list = append(list, u)
	}
	return list
}

// CreateUser creates a new user
func CreateUser(name, email string) User {
	mutex.Lock()
	defer mutex.Unlock()

	u := User{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
	}
	users[u.ID] = u
	return u
}

// GetUserByID returns a user by ID
func GetUserByID(id string) (User, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	u, ok := users[id]
	return u, ok
}

// ClearUsers clears all users (for testing purposes)
func ClearUsers() {
	mutex.Lock()
	defer mutex.Unlock()
	users = make(map[string]User)
}
