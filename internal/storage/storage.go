package storage

import (
	"github.com/outcatcher/anwil/internal/storage/inmemory"
	"github.com/outcatcher/anwil/internal/storage/schema"
)

var inMemDB = inmemory.NewDB() //nolint:gochecknoglobals

// Storage returning storage to be used.
func Storage() schema.Storage {
	return inMemDB
}
