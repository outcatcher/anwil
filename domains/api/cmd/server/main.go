/*
Package main is a script to start API server.
*/
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/outcatcher/anwil/domains/api"
	"github.com/outcatcher/anwil/domains/core/logging"
)

const defaultTimeout = time.Minute

func main() {
	logger := logging.GetDefaultLogger()

	argConfigPath := flag.String("config", "", "Configuration path")
	flag.Parse()

	if *argConfigPath == "" {
		logger.Fatalf("please provide configuration path")
	}

	logger.Printf("using configuration at %s", *argConfigPath)

	ctx := logging.CtxWithLogger(context.Background(), logger)

	if err := exec(ctx, *argConfigPath); err != nil {
		logger.Fatal(err)
	}
}

func exec(ctx context.Context, configPath string) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	state, err := api.Init(ctx, configPath)
	if err != nil {
		return fmt.Errorf("error initializing API: %w", err)
	}

	app, err := state.App() //nolint:contextcheck
	if err != nil {
		return fmt.Errorf("error serving HTTP: %w", err)
	}

	go func() {
		sig := <-sigChan

		logger := logging.LoggerFromCtx(ctx)

		logger.Printf("received signal: %+v", sig)

		err := app.ShutdownWithTimeout(defaultTimeout)
		if err != nil {
			logger.Printf("server shutdown faced error: %s", err)
		}
	}()

	cfg := state.Config()

	addr := fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port)

	loggedAddr := addr
	if cfg.API.Host == "" {
		loggedAddr = fmt.Sprintf("localhost:%d", cfg.API.Port)
	}

	state.Logger().Printf("Anwil API server started at http://%s", loggedAddr)

	err = app.Listen(addr)
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server stopped with error: %w", err)
	}

	return nil
}
