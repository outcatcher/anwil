/*
Package main is a script to run migrations.
*/
package main

import (
	"context"
	"flag"
	"log"
	"log/syslog"
	"path/filepath"

	"github.com/outcatcher/anwil/domains/internals/config"
	"github.com/outcatcher/anwil/domains/internals/storage"
)

func main() {
	sysLogger, err := syslog.New(syslog.LOG_INFO, "anwil-migrate")
	if err != nil {
		log.Fatalln(err)
	}

	log.SetOutput(sysLogger)

	argConfigPath := flag.String("config", "", "Configuration path")
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
		log.Fatal(err)
	}

	if err := storage.ApplyMigrations(cfg.DB); err != nil {
		log.Fatal(err)
	}
}
