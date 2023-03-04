package dto

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	errStateWithoutStorage   = errors.New("given state has no storage")
	errServiceWithoutStorage = errors.New("given service does not support storage")
)

// WithStorage defines service or state having storage attached.
type WithStorage interface {
	Storage() QueryExecutor
}

// RequiresStorage defines service which can use storage attached.
type RequiresStorage interface {
	UseStorage(executor QueryExecutor)
}

// InitWithStorage adds storage to the service.
func InitWithStorage(serv, state interface{}) error {
	reqStorage, ok := serv.(RequiresStorage)
	if !ok {
		return fmt.Errorf("error intializing service storage: %w", errServiceWithoutStorage)
	}

	stateWithStorage, ok := state.(WithStorage)
	if !ok {
		return fmt.Errorf("error intializing service storage: %w", errStateWithoutStorage)
	}

	reqStorage.UseStorage(stateWithStorage.Storage())

	return nil
}

// QueryExecutor interface describing sqlx.DB or sqlx.Tx in scope of the project.
type QueryExecutor interface {
	sqlx.ExtContext

	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}
