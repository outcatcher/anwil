/*
Package storage contains db-related operations with users.
*/
package storage

import (
	"context"
	"fmt"
)

// InsertUser creates a user.
func (u *UserStorage) InsertUser(ctx context.Context, data Wisher) error {
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
func (u *UserStorage) GetUser(ctx context.Context, username string) (*Wisher, error) {
	user := new(Wisher)

	err := u.db.GetContext(ctx, user, `SELECT * FROM wishers WHERE username = $1;`, username)
	if err != nil {
		return nil, fmt.Errorf("error selecting user: %w", err)
	}

	return user, nil
}
