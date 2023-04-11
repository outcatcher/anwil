package api

import (
	"context"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recov "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/outcatcher/anwil/domains/api/handlers"
	"github.com/outcatcher/anwil/domains/core/config"
	configSchema "github.com/outcatcher/anwil/domains/core/config/schema"
	"github.com/outcatcher/anwil/domains/core/logging"
	"github.com/outcatcher/anwil/domains/core/services"
	svcSchema "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/outcatcher/anwil/domains/storage"
	storageSchema "github.com/outcatcher/anwil/domains/storage/schema"
	users "github.com/outcatcher/anwil/domains/users/service"
	usersSchema "github.com/outcatcher/anwil/domains/users/service/schema"
	"github.com/valyala/fasthttp"
)

const defaultTimeout = time.Minute

// State holds general application state.
type State struct {
	cfg *configSchema.Configuration

	log     *log.Logger
	storage storageSchema.QueryExecutor

	services svcSchema.ServiceMapping
}

// App creates new API instance.
func (s *State) App() (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		StrictRouting:     false,
		CaseSensitive:     false,
		PassLocalsToViews: false,
		ReadTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
		IdleTimeout:       defaultTimeout,
		AppName:           "anwil",
	})

	app.Use(
		logger.New(logger.Config{Output: s.Logger().Writer()}),
		recov.New(recov.Config{EnableStackTrace: true}),
	)

	if err := handlers.PopulateEndpoints(app, s); err != nil { //nolint:contextcheck
		return nil, fmt.Errorf("error populating endpoints: %w", err)
	}

	return app, nil
}

// Server creates new API server instance.
func (s *State) Server() (*fasthttp.Server, error) {
	app, err := s.App()
	if err != nil {
		return nil, err
	}

	return app.Server(), nil
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
func (s *State) Logger() *log.Logger {
	return s.log
}

// Config returns server configuration.
func (s *State) Config() *configSchema.Configuration {
	return s.cfg
}

// Users service.
func (s *State) Users() usersSchema.Service {
	return s.services[usersSchema.ServiceUsers].(usersSchema.Service) //nolint:forcetypeassert
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
		log:     logging.LoggerFromCtx(ctx),
		storage: db,
	}

	apiState.WithServices(users.New())

	initialized, err := services.Initialize(ctx, apiState, apiState.services)
	if err != nil {
		return nil, fmt.Errorf("error initializing API: %w", err)
	}

	apiState.services = initialized

	return apiState, nil
}
