package inmemory

import (
	"context"
	"fmt"

	"github.com/outcatcher/anwil/internal/storage/schema"
)

// GetUser returns existing user by username.
func (p *db) GetUser(_ context.Context, username string) (*schema.User, error) {
	user, ok := p.users.Get(username)
	if !ok {
		return nil, fmt.Errorf("missing user %s: %w", username, schema.ErrNotFound)
	}

	return &user, nil
}

// SaveUser creates new or updates existing user.
func (p *db) SaveUser(_ context.Context, user schema.User) error {
	p.users.Set(user.Username, user)

	return nil
}
