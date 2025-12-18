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
	nextID = 1 // This will be removed once uuid is fully integrated
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

// CreateUser creates a new user
func CreateUser(u User) User {
	mutex.Lock()
	defer mutex.Unlock()
	u.ID = uuid.New().String() // Using uuid for ID generation
	users[u.ID] = u
	return u
}

// GetUserByID retrieves a user by ID
func GetUserByID(id string) (User, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	u, ok := users[id]
	return u, ok
}
