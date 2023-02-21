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
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/api/handlers"
	"github.com/outcatcher/anwil/config"
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

	if err := exec(configPath); err != nil {
		log.Fatal(err)
	}
}

func exec(configPath string) error {
	cfg, err := config.LoadServerConfiguration(path.Clean(configPath))
	if err != nil {
		return fmt.Errorf("error loading server config: %w", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	router, err := handlers.NewRouter(".", cfg, gin.Logger(), gin.Recovery())
	if err != nil {
		return fmt.Errorf("error creating new router: %w", err)
	}

	server := http.Server{ //nolint:exhaustruct
		Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: defaultTimeout,
	}

	log.Printf("Anwil API server started at http://%s", server.Addr)

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
