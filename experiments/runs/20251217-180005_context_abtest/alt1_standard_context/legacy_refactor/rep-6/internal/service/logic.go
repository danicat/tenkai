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
	users  = make(map[string]User)
	mutex  sync.RWMutex
)


// GetAllUsers returns a list of all users
func GetAllUsers() []User {
	mutex.RLock()
	defer mutex.RUnlock()
	var list []User
	for _, u := range users {
		list = append(list, u)
	}
	return list
}

// CreateUser adds a new user to the system
func CreateUser(u User) User {
	mutex.Lock()
	defer mutex.Unlock()
	
	// Original main.go used integer-based ID, but now using UUID as per dependency upgrade and best practices.
	u.ID = uuid.New().String()
	
	users[u.ID] = u
	return u
}

// GetUserByID retrieves a user by their ID
func GetUserByID(id string) (User, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	u, ok := users[id]
	return u, ok
}

// ResetUsersForTest clears the user store for testing purposes
func ResetUsersForTest() {
	mutex.Lock()
	defer mutex.Unlock()
	users = make(map[string]User)
}
