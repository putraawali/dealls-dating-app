package usecase_helpers

import (
	"dealls-dating-app/src/models"
	"sync"
)

type temporaryAvailableUser struct {
	mu sync.Mutex

	// Store available user
	availableUser []models.User

	// To map inserted user, if already inserted, will not append to availableUser
	insertedUser map[int64]bool
}

func NewHandleAvailableUser() temporaryAvailableUser {
	return temporaryAvailableUser{
		insertedUser:  make(map[int64]bool),
		availableUser: []models.User{},
	}
}

func (t *temporaryAvailableUser) Add(data models.User) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.insertedUser[data.UserID] {
		t.availableUser = append(t.availableUser, data)
		t.insertedUser[data.UserID] = true
	}
}

func (t *temporaryAvailableUser) GetCurrent() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	return len(t.availableUser)
}

func (t *temporaryAvailableUser) FinalData() []models.User {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.availableUser
}
