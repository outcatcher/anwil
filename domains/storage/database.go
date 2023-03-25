package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // use `postgres` driver
	"github.com/outcatcher/anwil/domains/core/config/schema"
)

const dbDriver = "postgres"

func dbString(dbConfig schema.DatabaseConfiguration) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.Username, dbConfig.Password,
		dbConfig.Host, dbConfig.Port,
		dbConfig.DatabaseName,
	)
}

// Connect connects to the database with given configuration.
func Connect(cfg schema.DatabaseConfiguration) (*sqlx.DB, error) {
	db, err := sqlx.Connect(dbDriver, dbString(cfg))
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}
