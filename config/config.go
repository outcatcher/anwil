/*
Package config handles application configuration loading.
*/
package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// ServerConfiguration - server configuration.
type ServerConfiguration struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	KeyStoragePath string `yaml:"key_storage_path"`
}

// LoadServerConfiguration loads server yaml configuration by given path.
// It strictly validates yaml file contents, so will fail in case yaml structure is incorrect.
func LoadServerConfiguration(path string) (*ServerConfiguration, error) {
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

	decoder := yaml.NewDecoder(file)
	decoder.SetStrict(true)

	cfg := new(ServerConfiguration)

	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("config decode error: %w", err)
	}

	return cfg, nil
}
