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
	Users = make(map[string]User)
	Mutex sync.RWMutex
)

// GetAllUsers returns a list of all users
func GetAllUsers() []User {
	Mutex.RLock()
	defer Mutex.RUnlock()
	var list []User
	for _, u := range Users {
		list = append(list, u)
	}
	return list
}

// CreateUser creates a new user with a generated ID
func CreateUser(u User) (User, error) {
	Mutex.Lock()
	defer Mutex.Unlock()
	u.ID = uuid.New().String()
	Users[u.ID] = u
	return u, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(id string) (User, bool) {
	Mutex.RLock()
	defer Mutex.RUnlock()
	u, ok := Users[id]
	return u, ok
}
