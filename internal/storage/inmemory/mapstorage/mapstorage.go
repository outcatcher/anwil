/*
Package mapstorage contains implementation of map-based in-memory storage.
*/
package mapstorage

import (
	"math/rand"
	"sync"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixMicro())) //nolint:gosec

// Storage is setting, storing and returning values by keys.
type Storage[kT comparable, vT any] interface {
	Get(key kT) (vT, bool)
	Set(key kT, value vT)
}

// mapStorage is map-based in-memory storage with lock.
type mapStorage[kT comparable, vT any] struct {
	storage map[kT]vT

	lock sync.Mutex
}

// Get returns storage[key].
func (m *mapStorage[kT, vT]) Get(key kT) (vT, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.storage == nil {
		m.storage = make(map[kT]vT)
	}

	value, ok := m.storage[key]

	return value, ok
}

// Set sets storage[key] = value.
func (m *mapStorage[kT, vT]) Set(key kT, value vT) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.storage == nil {
		m.storage = make(map[kT]vT)
	}

	m.storage[key] = value
}

// NewStorage creates new map storage.
func NewStorage[kT comparable, vT any]() Storage[kT, vT] {
	return &mapStorage[kT, vT]{
		storage: make(map[kT]vT),
		lock:    sync.Mutex{},
	}
}

// GenerateRandomInt64ID generates random ID not existing in the storage.
func GenerateRandomInt64ID[vT any](storage Storage[int64, vT]) int64 {
	value := random.Int63()

	for { // for mega-luck champions hitting already generated int64 numbers
		if _, ok := storage.Get(value); !ok {
			break
		}

		value = random.Int63()
	}

	return value
}
