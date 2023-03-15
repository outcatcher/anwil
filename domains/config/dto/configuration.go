/*
Package dto contains shared configuration definitions.

There must be no use of other DTOs in this package.
*/
package dto

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/outcatcher/anwil/domains/logging"
)

const encodedPrivateKeyLen = 128

var (
	errStateWithoutConfig   = errors.New("given state has no config")
	errServiceWithoutConfig = errors.New("given service does not support config")
)

// WithConfig can return configuration.
type WithConfig interface {
	// Config returns used configuration.
	Config() *Configuration
}

// RequiresConfig can use configuration.
type RequiresConfig interface {
	// UseConfig attaches configuration to the service.
	UseConfig(*Configuration)
}

// ConfigInject injects configuration into service.
func ConfigInject(service interface{}, state interface{}) error {
	reqConfig, ok := service.(RequiresConfig)
	if !ok {
		return fmt.Errorf("error intializing service config: %w", errServiceWithoutConfig)
	}

	stateWithConfig, ok := state.(WithConfig)
	if !ok {
		return fmt.Errorf("error intializing service config: %w", errStateWithoutConfig)
	}

	reqConfig.UseConfig(stateWithConfig.Config())

	return nil
}

// DatabaseConfiguration - DB-related configuration.
//
// Note that for fields with `env` tag, environment variable value has priority over yaml value.
type DatabaseConfiguration struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	DatabaseName  string `yaml:"databaseName"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password" env:"POSTGRES_PASSWORD"`
	MigrationsDir string `yaml:"migrationsDir"`
}

// APIConfiguration - API-related configuration.
type APIConfiguration struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	StaticPath string `yaml:"staticPath"`
}

// Configuration - overall system configuration.
type Configuration struct {
	API            APIConfiguration      `yaml:"api"`
	DB             DatabaseConfiguration `yaml:"db"`
	PrivateKeyPath string                `yaml:"privateKeyPath"`
	Debug          bool                  `yaml:"debug"`

	privateKey ed25519.PrivateKey
}

// loadPrivateKey loads ED25519 private key from path.
func loadPrivateKey(ctx context.Context, privateKeyPath string) (ed25519.PrivateKey, error) {
	keyFile, err := os.Open(filepath.Clean(privateKeyPath))
	if err != nil {
		return nil, fmt.Errorf("error opening private key file: %w", err)
	}

	defer func() {
		closeErr := keyFile.Close()
		if closeErr != nil {
			logging.LoggerFromCtx(ctx).Println(closeErr)
		}
	}()

	// reading with fixed size allows file to have \n or \r\n after key data
	keyData := make([]byte, encodedPrivateKeyLen)

	_, err = io.ReadFull(keyFile, keyData)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %w", err)
	}

	decodedData := make([]byte, hex.DecodedLen(len(keyData)))
	if _, err := hex.Decode(decodedData, keyData); err != nil {
		return nil, fmt.Errorf("error decoding loaded key file: %w", err)
	}

	return decodedData, nil
}

// GetPrivateKey returns private key value from the loaded configuration.
func (s Configuration) GetPrivateKey(ctx context.Context) (ed25519.PrivateKey, error) {
	if len(s.privateKey) > 0 {
		return s.privateKey, nil
	}

	privateKey, err := loadPrivateKey(ctx, s.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error getting loading private key: %w", err)
	}

	s.privateKey = privateKey

	return s.privateKey, nil
}
