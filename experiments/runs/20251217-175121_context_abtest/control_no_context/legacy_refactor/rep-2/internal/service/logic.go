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
	// nextID is no longer needed as we'll use UUIDs
)

func GetAllUsers() []User {
	mutex.RLock()
	defer mutex.RUnlock()
	var list []User
	for _, u := range users {
		list = append(list, u)
	}
	return list
}

func CreateUser(user User) User {
	mutex.Lock()
	defer mutex.Unlock()
	user.ID = uuid.New().String() // Use UUID for ID generation
	users[user.ID] = user
	return user
}

func GetUserByID(id string) (User, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	user, ok := users[id]
	return user, ok
}