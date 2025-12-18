package service

import (
	"sync"

	"github.com/google/uuid"
	"rep-1/internal/models"
)

// In-memory storage for users
var (
	users = make(map[string]models.User)
	mutex sync.RWMutex
)

// GetAllUsers returns a list of all users
func GetAllUsers() []models.User {
	mutex.RLock()
	defer mutex.RUnlock()
	var list []models.User
	for _, u := range users {
		list = append(list, u)
	}
	return list
}

// CreateUser creates a new user and adds it to the storage
func CreateUser(user models.User) (models.User, error) {
	mutex.Lock()
	defer mutex.Unlock()

	user.ID = uuid.New().String()
	users[user.ID] = user
	return user, nil
}

// GetUserByID returns a user by their ID
func GetUserByID(id string) (models.User, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	user, ok := users[id]
	return user, ok
}
