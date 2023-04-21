/*
Package schema contains storage-related DTOs
*/
package schema

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/outcatcher/anwil/domains/core/services"
)

// WithStorage defines service or state having storage attached.
type WithStorage interface {
	Storage() QueryExecutor
}

// RequiresStorage defines service which can use storage attached.
type RequiresStorage interface {
	UseStorage(executor QueryExecutor)
}

// StorageInject adds storage to the service.
func StorageInject(consumer, provider any) error {
	reqStorage, provStorage, err := services.ValidateArgInterfaces[RequiresStorage, WithStorage](consumer, provider)
	if err != nil {
		return fmt.Errorf("error injecting storage: %w", err)
	}

	reqStorage.UseStorage(provStorage.Storage())

	return nil
}

// QueryExecutor interface describing sqlx.DB or sqlx.Tx in scope of the project.
type QueryExecutor interface {
	sqlx.ExtContext

	GetContext(ctx context.Context, dest any, query string, args ...any) error
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
}
