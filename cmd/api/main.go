/*
Package main is a script to start API server.
*/
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/outcatcher/anwil/domains/api"
)

const defaultTimeout = time.Minute

func main() {
	argConfigPath := flag.String("config", "", "Configuration path")
	flag.Parse()

	if *argConfigPath == "" {
		log.Fatalf("please provide configuation path")
	}

	configPath, err := filepath.Abs(filepath.Clean(*argConfigPath))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("using configuration at %s", configPath)

	if err := exec(context.Background(), configPath); err != nil {
		log.Fatal(err)
	}
}

func exec(ctx context.Context, configPath string) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	state, err := api.Init(ctx, configPath)
	if err != nil {
		return fmt.Errorf("error initializing API: %w", err)
	}

	server, err := state.Server(ctx)
	if err != nil {
		return fmt.Errorf("error serving HTTP: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	go func() {
		sig := <-sigChan

		log.Printf("received signal: %+v", sig)

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("server shutdown faced error: %s", err)
		}
	}()

	err = server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server stopped with error: %w", err)
	}

	return nil
}
