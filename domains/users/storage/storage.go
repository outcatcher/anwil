package storage

import (
	"context"
)

// UserStorage - хранилище данных пользователя.
type UserStorage interface {
	InsertUser(ctx context.Context, data Wisher) error
	GetUser(ctx context.Context, username string) (*Wisher, error)
}
