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

// CreateUser creates a new user and returns it
func CreateUser(user User) User {
	mutex.Lock()
	defer mutex.Unlock()
	// Using uuid.New().String() as per the plan to upgrade uuid
	user.ID = uuid.New().String()
	users[user.ID] = user
	return user
}

// GetUserByID returns a user by ID
func GetUserByID(id string) (User, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	u, ok := users[id]
	return u, ok
}
