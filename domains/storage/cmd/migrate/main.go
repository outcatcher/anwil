/*
Package main is a script to run migrations.
*/
package main

import (
	"context"
	"flag"
	"log"
	"path/filepath"

	"github.com/outcatcher/anwil/domains/core/config"
	"github.com/outcatcher/anwil/domains/core/logging"
	"github.com/outcatcher/anwil/domains/storage"
)

func main() {
	logger := logging.GetDefaultLogger()

	log.SetOutput(logger.Writer()) // for subsequent goose calls

	argConfigPath := flag.String("config", "", "Configuration path")
	flag.Parse()

	if *argConfigPath == "" {
		logger.Fatalf("please provide configuation path")
	}

	configPath, err := filepath.Abs(filepath.Clean(*argConfigPath))
	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()

	cfg, err := config.LoadServerConfiguration(ctx, configPath)
	if err != nil {
		logger.Fatal(err)
	}

	if err := storage.ApplyMigrations(cfg.DB); err != nil {
		logger.Fatal(err)
	}
}
