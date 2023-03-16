package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/api/handlers"
	"github.com/outcatcher/anwil/domains/internals/config"
	configSchema "github.com/outcatcher/anwil/domains/internals/config/schema"
	"github.com/outcatcher/anwil/domains/internals/logging"
	"github.com/outcatcher/anwil/domains/internals/services"
	svcSchema "github.com/outcatcher/anwil/domains/internals/services/schema"
	"github.com/outcatcher/anwil/domains/internals/storage"
	storageDTO "github.com/outcatcher/anwil/domains/internals/storage/schema"
	users "github.com/outcatcher/anwil/domains/users/service"
	usersSchema "github.com/outcatcher/anwil/domains/users/service/schema"
)

const defaultTimeout = time.Minute

// State holds general application state.
type State struct {
	cfg *configSchema.Configuration

	log     *log.Logger
	storage storageDTO.QueryExecutor

	services svcSchema.ServiceMapping
}

// Server creates new API server instance.
func (s *State) Server(ctx context.Context) (*http.Server, error) {
	cfg := s.Config()

	engine := gin.New()
	engine.Use(gin.LoggerWithWriter(s.Logger().Writer()), gin.Recovery())

	// запросы не должны использовать родительский контекст
	if err := handlers.PopulateEndpoints(engine, s); err != nil { //nolint:contextcheck
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

// WithServices uses selected service mapping.
func (s *State) WithServices(services ...svcSchema.Service) {
	if s.services == nil {
		s.services = make(svcSchema.ServiceMapping)
	}

	for _, svc := range services {
		s.services[svc.ID()] = svc
	}
}

// Logger returns configured logger or a default one.
func (s *State) Logger() *log.Logger {
	if s.log == nil {
		s.log = log.Default()
	}

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
func (s *State) Storage() storageDTO.QueryExecutor {
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
