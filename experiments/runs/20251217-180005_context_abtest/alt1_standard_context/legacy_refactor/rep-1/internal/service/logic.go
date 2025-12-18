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

// CreateUser creates a new user and adds it to the system
func CreateUser(name, email string) User {
	mutex.Lock()
	defer mutex.Unlock()

	newUser := User{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
	}
	users[newUser.ID] = newUser
	return newUser
}

// GetUser retrieves a user by ID
func GetUser(id string) (User, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	user, ok := users[id]
	return user, ok
}

// GetAllUsers retrieves all users in the system
func GetAllUsers() []User {
	mutex.RLock()
	defer mutex.RUnlock()
	var list []User
	for _, u := range users {
		list = append(list, u)
	}
	return list
}
