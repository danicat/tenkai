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
	nextID = 1
)

// GetUsers returns a list of all users
func GetUsers() []User {
	mutex.RLock()
	defer mutex.RUnlock()
	var list []User
	for _, u := range users {
		list = append(list, u)
	}
	return list
}

// GetUser returns a user by ID
func GetUser(id string) (User, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	u, ok := users[id]
	return u, ok
}

// CreateUser creates a new user and returns it
func CreateUser(u User) User {
	mutex.Lock()
	defer mutex.Unlock()
	u.ID = uuid.New().String()
	users[u.ID] = u
	return u
}
