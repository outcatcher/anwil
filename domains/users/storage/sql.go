/*
Package storage contains db-related operations with users.
*/
package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	services "github.com/outcatcher/anwil/domains/internals/services/schema"
	storageSchema "github.com/outcatcher/anwil/domains/internals/storage/schema"
)

// userStorage - storage of users.
type userStorage struct {
	db storageSchema.QueryExecutor
}

// New create a new UserStorage instance.
func New(db storageSchema.QueryExecutor) UserStorage {
	return &userStorage{db: db}
}

// InsertUser creates a user.
func (u *userStorage) InsertUser(ctx context.Context, data Wisher) error {
	_, err := u.db.NamedExecContext(
		ctx,
		`INSERT INTO wishers (username, password, full_name) VALUES (:username, :password, :full_name);`,
		data,
	)
	if err != nil {
		return fmt.Errorf("inserting user failed: %w", err)
	}

	return nil
}

// GetUser returns single user by username.
func (u *userStorage) GetUser(ctx context.Context, username string) (*Wisher, error) {
	user := new(Wisher)

	err := u.db.GetContext(ctx, user, `SELECT * FROM wishers WHERE username = $1;`, username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no user found: %w", services.ErrNotFound)
	}

	if err != nil {
		return nil, fmt.Errorf("error selecting user: %w", err)
	}

	return user, nil
}
