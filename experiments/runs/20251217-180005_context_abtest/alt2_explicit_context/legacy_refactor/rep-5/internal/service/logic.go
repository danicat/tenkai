package service

import (
	"fmt"
	"sync"
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
	nextID = 1
)

// InitService resets the service's state for testing
func InitService() {
	mutex.Lock()
	defer mutex.Unlock()
	users = make(map[string]User)
	nextID = 1
}

// GetAllUsers returns a list of all users
func GetAllUsers() []User {
	mutex.RLock()
	defer mutex.RUnlock()
	var list = make([]User, 0)
	for _, u := range users {
		list = append(list, u)
	}
	return list
}

// CreateUser creates a new user
func CreateUser(u User) User {
	mutex.Lock()
	defer mutex.Unlock()
	u.ID = fmt.Sprintf("%d", nextID)
	nextID++
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
