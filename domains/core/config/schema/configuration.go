/*
Package schema contains shared configuration definitions.

There must be no use of other DTOs in this package.
*/
package schema

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/outcatcher/anwil/domains/core/services"
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
func ConfigInject(consumer, provider any) error {
	reqConfig, provConfig, err := services.ValidateArgInterfaces[RequiresConfig, WithConfig](consumer, provider)
	if err != nil {
		return fmt.Errorf("error injecting configuration: %w", err)
	}

	reqConfig.UseConfig(provConfig.Config())

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

	decodedData := make([]byte, ed25519.PrivateKeySize)

	_, err = io.ReadFull(hex.NewDecoder(keyFile), decodedData)
	if err != nil {
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
