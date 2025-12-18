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

func CreateUser(name, email string) User {
	mutex.Lock()
	defer mutex.Unlock()

	uuid := uuid.New().String()
	u := User{
		ID:    uuid,
		Name:  name,
		Email: email,
	}
	users[u.ID] = u
	return u
}

func GetUsers() []User {
	mutex.RLock()
	defer mutex.RUnlock()

	var list []User
	for _, u := range users {
		list = append(list, u)
	}
	return list
}

func GetUserByID(id string) (User, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	u, ok := users[id]
	return u, ok
}
