/*
Package config handles application configuration loading.
*/
package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/outcatcher/anwil/domains/config/dto"
	"github.com/outcatcher/anwil/domains/config/envyaml"
)

// LoadServerConfiguration loads server yaml configuration by given path.
// It strictly validates yaml file contents, so will fail in case yaml structure is incorrect.
func LoadServerConfiguration(ctx context.Context, path string) (*dto.Configuration, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path from %s: %w", path, err)
	}

	file, err := os.Open(filepath.Clean(absPath))
	if err != nil {
		return nil, fmt.Errorf("error opening configuration file %s: %w", absPath, err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("error closing config file: %v", err)
		}
	}()

	cfg := new(dto.Configuration)

	if err := envyaml.Decode(ctx, file, cfg); err != nil {
		return nil, fmt.Errorf("config decode error: %w", err)
	}

	return cfg, nil
}
