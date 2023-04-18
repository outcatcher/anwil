package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/outcatcher/anwil/domains/api/errorhandler"
	"github.com/outcatcher/anwil/domains/api/handlers"
	"github.com/outcatcher/anwil/domains/api/middlewares"
	"github.com/outcatcher/anwil/domains/core/config"
	configSchema "github.com/outcatcher/anwil/domains/core/config/schema"
	"github.com/outcatcher/anwil/domains/core/services"
	svcSchema "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/outcatcher/anwil/domains/storage"
	storageSchema "github.com/outcatcher/anwil/domains/storage/schema"
	users "github.com/outcatcher/anwil/domains/users/service"
)

const defaultTimeout = time.Minute

// State holds general application state.
type State struct {
	cfg *configSchema.Configuration

	storage storageSchema.QueryExecutor

	services svcSchema.ServiceMapping
}

// Server creates new API server instance.
func (s *State) Server(ctx context.Context) (*http.Server, error) {
	cfg := s.Config()

	engine := echo.New()

	engine.HTTPErrorHandler = errorhandler.HandleErrors()

	engine.Use(
		middleware.LoggerWithConfig(middleware.LoggerConfig{Output: s.Logger().Writer()}),
		middleware.Recover(),
		middleware.RemoveTrailingSlash(),
		middlewares.RequireJSON,
	)

	// запросы не должны использовать родительский контекст
	if err := handlers.PopulateEndpoints(ctx, engine, s); err != nil {
		return nil, fmt.Errorf("error populating endpoints: %w", err)
	}

	server := &http.Server{ //nolint:exhaustruct
		Addr:              fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port),
		Handler:           engine,
		ReadHeaderTimeout: defaultTimeout,
		BaseContext:       func(_ net.Listener) context.Context { return ctx },
	}

	loggedAddr := server.Addr

	if cfg.API.Host == "" {
		loggedAddr = fmt.Sprintf("localhost:%d", cfg.API.Port)
	}

	s.Logger().Printf("Anwil API server started at http://%s", loggedAddr)

	return server, nil
}

// WithServices uses selected services.
func (s *State) WithServices(services ...svcSchema.Service) {
	if s.services == nil {
		s.services = make(svcSchema.ServiceMapping)
	}

	for _, svc := range services {
		s.services[svc.ID()] = svc
	}
}

// Logger returns configured logger.
func (*State) Logger() *log.Logger {
	return log.Default()
}

// Config returns server configuration.
func (s *State) Config() *configSchema.Configuration {
	return s.cfg
}

// Service returns exact service instance by ID.
func (s *State) Service(id svcSchema.ServiceID) svcSchema.Service {
	return s.services[id]
}

// Storage returns shared query executor (i.e. *sqlx.DB).
func (s *State) Storage() storageSchema.QueryExecutor {
	return s.storage
}

// Init initializes API and returns new API instance.
func Init(ctx context.Context, configPath string) (*State, error) {
	cfg, err := config.LoadServerConfiguration(ctx, path.Clean(configPath))
	if err != nil {
		return nil, fmt.Errorf("error loading server config: %w", err)
	}

	db, err := storage.Connect(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the storage: %w", err)
	}

	apiState := &State{
		cfg:     cfg,
		storage: db,
	}

	apiState.WithServices(new(users.Service))

	initialized, err := services.Initialize(ctx, apiState, apiState.services)
	if err != nil {
		return nil, fmt.Errorf("error initializing API: %w", err)
	}

	apiState.services = initialized

	return apiState, nil
}
