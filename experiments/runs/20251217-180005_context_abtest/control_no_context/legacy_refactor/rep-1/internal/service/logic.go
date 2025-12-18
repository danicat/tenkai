package service

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

var (
	users = make(map[uuid.UUID]User)
	mutex sync.RWMutex
)

// ClearUsers clears the users map (for testing purposes)
func ClearUsers() {
	mutex.Lock()
	defer mutex.Unlock()
	users = make(map[uuid.UUID]User)
}

// CreateUser creates a new user and adds it to the system
func CreateUser(name, email string) (User, error) {
	mutex.Lock()
	defer mutex.Unlock()

	id, err := uuid.NewRandom()
	if err != nil {
		return User{}, fmt.Errorf("failed to generate UUID: %w", err)
	}

	user := User{
		ID:    id,
		Name:  name,
		Email: email,
	}
	users[user.ID] = user
	return user, nil
}

// GetAllUsers returns a list of all users
func GetAllUsers() []User {
	mutex.RLock()
	defer mutex.RUnlock()

	var userList []User
	for _, u := range users {
		userList = append(userList, u)
	}
	return userList
}

// GetUserByID returns a user by their ID
func GetUserByID(id uuid.UUID) (User, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	user, ok := users[id]
	return user, ok
}
