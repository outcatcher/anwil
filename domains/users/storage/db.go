package storage

import (
	storageDTO "github.com/outcatcher/anwil/domains/storage/dto"
)

// UserStorage - storage of users.
type UserStorage struct {
	db storageDTO.QueryExecutor
}

// New create a new UserStorage instance.
func New(db storageDTO.QueryExecutor) *UserStorage {
	return &UserStorage{db: db}
}
