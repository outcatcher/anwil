package inmemory

import (
	"github.com/outcatcher/anwil/internal/storage/inmemory/mapstorage"
	"github.com/outcatcher/anwil/internal/storage/schema"
)

// db - in-memory implementation of Storage.
type db struct {
	users     mapstorage.Storage[string, schema.User]
	wishlists mapstorage.Storage[int64, schema.Wishlist]
}

// NewDB returns pre-initialized storage.
func NewDB() schema.Storage {
	return &db{
		users:     mapstorage.NewStorage[string, schema.User](),
		wishlists: mapstorage.NewStorage[int64, schema.Wishlist](),
	}
}
