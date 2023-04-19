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
	log.SetOutput(logging.GetDefaultLogWriter()) // for subsequent goose calls

	argConfigPath := flag.String("config", "", "Configuration path")
	argCommand := flag.String("command", "up", "Command for goose to execute")

	flag.Parse()

	if *argConfigPath == "" {
		log.Fatalf("please provide configuation path")
	}

	configPath, err := filepath.Abs(filepath.Clean(*argConfigPath))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	cfg, err := config.LoadServerConfiguration(ctx, configPath)
	if err != nil {
		log.Fatalln(err)
	}

	if err := storage.ApplyMigrations(cfg.DB, *argCommand); err != nil {
		log.Fatal(err)
	}
}
