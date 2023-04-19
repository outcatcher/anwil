package storage

import (
	"fmt"
	"path/filepath"

	config "github.com/outcatcher/anwil/domains/core/config/schema"
	"github.com/pressly/goose/v3"
)

// ApplyMigrations applies all available migrations.
func ApplyMigrations(cfg config.DatabaseConfiguration, command string) error {
	if err := goose.SetDialect(dbDriver); err != nil {
		return fmt.Errorf("error selecting dialect for migrations: %w", err)
	}

	db, err := Connect(cfg)
	if err != nil {
		return err
	}

	// we can't use dbDriver here - different drivers can be used for same DB engine
	migrationsPath := filepath.Clean(filepath.Join(cfg.MigrationsDir, "postgres"))

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("error getting abs path for %s: %w", migrationsPath, err)
	}

	if err := goose.RunWithOptions(command, db.DB, absPath, nil); err != nil {
		return fmt.Errorf("error applying migrations: %w", err)
	}

	return nil
}
