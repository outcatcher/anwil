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
	"github.com/outcatcher/anwil/domains/api/commonhandlers"
	"github.com/outcatcher/anwil/domains/api/errorhandler"
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
	// Shared configuration
	cfg *configSchema.Configuration

	// Shared storage driver, i.e. *sqlx.DB
	storage storageSchema.QueryExecutor

	// Actual initialized services
	services svcSchema.ServiceMapping
	// Functions to add handlers after HTTP server is created
	addHandlerFuncs []svcSchema.AddHandlersFunc
}

func (s *State) initEngine() (*echo.Echo, error) {
	engine := echo.New()

	engine.HTTPErrorHandler = errorhandler.HandleErrors()

	engine.Use(
		middleware.LoggerWithConfig(middleware.LoggerConfig{Output: s.Logger().Writer()}),
		middleware.Recover(),
		middleware.RemoveTrailingSlash(),
		middlewares.RequireJSON,
	)

	engine.Static("/static", s.Config().API.StaticPath)

	baseGroup := engine.Group("/api/v1")
	secGroup := baseGroup.Group("", middlewares.JWTAuth(s))

	for _, addHandlersFunc := range s.addHandlerFuncs {
		err := addHandlersFunc(baseGroup, secGroup)
		if err != nil {
			return nil, fmt.Errorf("error adding handlers for the service: %w", err)
		}
	}

	return engine, nil
}

// Server creates new API server instance.
func (s *State) Server(ctx context.Context) (*http.Server, error) {
	cfg := s.Config()

	engine, err := s.initEngine()
	if err != nil {
		return nil, fmt.Errorf("error creating server: %w", err)
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

// Logger returns configured logger.
func (*State) Logger() *log.Logger {
	return log.Default()
}

// Config returns server configuration.
func (s *State) Config() *configSchema.Configuration {
	return s.cfg
}

// Service returns exact service instance by ID.
func (s *State) Service(id svcSchema.ServiceID) any {
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

	usedServices := []svcSchema.ServiceDefinition{users.NewUserService()}

	initialized, err := services.Initialize(ctx, apiState, usedServices...)
	if err != nil {
		return nil, fmt.Errorf("error initializing API: %w", err)
	}

	apiState.services = initialized

	// Pre-initializing handlers. At this point there is no server, so populating the functions to be
	// called to add handlers when it will be ready.
	//
	// To be additionally considered: is there a need for service initialization before creating the server?
	addHandlerFuncs := make([]svcSchema.AddHandlersFunc, 1, len(usedServices)+1)

	addHandlerFuncs[0] = commonhandlers.AddEchoHandlers

	for _, svc := range usedServices {
		if svc.InitHandlersFunc == nil {
			continue // some services can possibly have no endpoints
		}

		addHandlerFuncs = append(addHandlerFuncs, svc.InitHandlersFunc(apiState))
	}

	apiState.addHandlerFuncs = addHandlerFuncs

	return apiState, nil
}
