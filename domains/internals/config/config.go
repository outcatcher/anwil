/*
Package config handles application configuration loading.
*/
package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/imdario/mergo"
	"github.com/outcatcher/anwil/domains/internals/config/schema"
	"github.com/outcatcher/anwil/domains/internals/logging"
	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v3"
)

const configFileCacheSize = 1 // now we have only one config

var configFileLRU, _ = lru.New[string, schema.Configuration](configFileCacheSize) // error is returned only on negative size

// LoadServerConfiguration loads server yaml configuration by given path and merges with defined env vars.
//
// It strictly validates yaml file contents, so will fail in case yaml structure is incorrect.
//
// File contents are hashed. Env vars are loaded each time function is called.
func LoadServerConfiguration(ctx context.Context, path string) (*schema.Configuration, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path from %s: %w", path, err)
	}

	value, ok := configFileLRU.Get(absPath)
	if ok {
		return &value, nil
	}

	file, err := os.Open(filepath.Clean(absPath))
	if err != nil {
		return nil, fmt.Errorf("error opening configuration file %s: %w", absPath, err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			logger := logging.LoggerFromCtx(ctx)

			logger.Printf("error closing config file: %v", err)
		}
	}()

	cfg := new(schema.Configuration)

	err = yaml.NewDecoder(file).Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("error decoding configuration file: %w", err)
	}

	fromEmv := new(schema.Configuration)

	if err := envconfig.Process(ctx, fromEmv); err != nil {
		return nil, fmt.Errorf("error loading config from env: %w", err)
	}

	if err := mergo.Merge(cfg, fromEmv, mergo.WithOverride); err != nil {
		return nil, fmt.Errorf("error merging configurations: %w", err)
	}

	configFileLRU.Add(absPath, *cfg)

	return cfg, nil
}
