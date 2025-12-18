package service

import (
	"fmt"
	"sync"

	"github.com/google/uuid" // Import the uuid package
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
	// nextID = 1 // No longer needed
)

// CreateUser creates a new user
func CreateUser(name, email string) (*User, error) {
	mutex.Lock()
	defer mutex.Unlock()

	u := User{
		ID:    uuid.New().String(), // Use UUID for ID generation
		Name:  name,
		Email: email,
	}
	// nextID++ // No longer needed
	users[u.ID] = u
	return &u, nil
}

// ListUsers returns all users
func ListUsers() []User {
	mutex.RLock()
	defer mutex.RUnlock()

	var list []User
	for _, u := range users {
		list = append(list, u)
	}
	return list
}

// GetUser retrieves a user by ID
func GetUser(id string) (*User, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	u, ok := users[id]
	if !ok {
		return nil, fmt.Errorf("user with ID %s not found", id)
	}
	return &u, nil
}

// ResetUsers clears the users map (for testing purposes)
func ResetUsers() {
	mutex.Lock()
	defer mutex.Unlock()
	users = make(map[string]User)
}
