/*
Package dto contains shared configuration definitions.

There must be no use of other DTOs in this package.
*/
package dto

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	errStateWithoutConfig   = errors.New("given state has no config")
	errServiceWithoutConfig = errors.New("given service does not support config")
)

// WithConfig can return configuration.
type WithConfig interface {
	Config() *Configuration
}

type RequiresConfig interface {
	UseConfig(*Configuration)
}

// InitWithConfig initializes given service with authentication.
func InitWithConfig(service interface{}, state interface{}) error {
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
func loadPrivateKey(privateKeyPath string) (ed25519.PrivateKey, error) {
	keyFile, err := os.Open(filepath.Clean(privateKeyPath))
	if err != nil {
		return nil, fmt.Errorf("error opening private key file: %w", err)
	}

	defer func() {
		closeErr := keyFile.Close()
		if closeErr != nil {
			log.Println(closeErr)
		}
	}()

	keyData, err := io.ReadAll(keyFile)
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
func (s Configuration) GetPrivateKey() (ed25519.PrivateKey, error) {
	if len(s.privateKey) > 0 {
		return s.privateKey, nil
	}

	privateKey, err := loadPrivateKey(s.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error getting loading private key: %w", err)
	}

	s.privateKey = privateKey

	return s.privateKey, nil
}
