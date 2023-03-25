package storage

import (
	"context"
	"fmt"
	"sync"

	services "github.com/outcatcher/anwil/domains/core/services/schema"
)

type mock struct {
	storage     map[string]Wisher
	storageLock sync.Mutex
}

// InsertUser inserts user into storage.
func (m *mock) InsertUser(_ context.Context, data Wisher) error {
	m.storageLock.Lock()
	defer m.storageLock.Unlock()

	if _, ok := m.storage[data.Username]; ok {
		return fmt.Errorf("username %s: %w", data.Username, services.ErrConflict)
	}

	m.storage[data.Username] = data

	return nil
}

// GetUser returns user from storage.
func (m *mock) GetUser(_ context.Context, username string) (*Wisher, error) {
	m.storageLock.Lock()
	defer m.storageLock.Unlock()

	user, ok := m.storage[username]
	if !ok {
		return nil, fmt.Errorf("username %s: %w", username, services.ErrNotFound)
	}

	return &user, nil
}

// NewMock create a new UserStorage mock instance.
func NewMock() UserStorage {
	return &mock{storage: make(map[string]Wisher)}
}
