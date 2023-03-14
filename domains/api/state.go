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
	"github.com/outcatcher/anwil/domains/auth"
	authDTO "github.com/outcatcher/anwil/domains/auth/dto"
	"github.com/outcatcher/anwil/domains/config"
	configDTO "github.com/outcatcher/anwil/domains/config/dto"
	"github.com/outcatcher/anwil/domains/logging"
	"github.com/outcatcher/anwil/domains/services"
	svcDTO "github.com/outcatcher/anwil/domains/services/dto"
	"github.com/outcatcher/anwil/domains/storage"
	storageDTO "github.com/outcatcher/anwil/domains/storage/dto"
	"github.com/outcatcher/anwil/domains/users"
	usersDTO "github.com/outcatcher/anwil/domains/users/dto"
)

const defaultTimeout = time.Minute

// State holds general application state.
type State struct {
	cfg *configDTO.Configuration

	log     *log.Logger
	storage storageDTO.QueryExecutor

	services svcDTO.ServiceMapping
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
func (s *State) WithServices(services ...svcDTO.Service) {
	if s.services == nil {
		s.services = make(svcDTO.ServiceMapping)
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
func (s *State) Config() *configDTO.Configuration {
	return s.cfg
}

// Authentication service.
func (s *State) Authentication() authDTO.Service {
	return s.services[authDTO.ServiceAuth].(authDTO.Service) //nolint:forcetypeassert
}

// Users service.
func (s *State) Users() usersDTO.Service {
	return s.services[usersDTO.ServiceUsers].(usersDTO.Service) //nolint:forcetypeassert
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

	apiState.WithServices(auth.New(), users.New())

	initialized, err := services.Initialize(ctx, apiState, apiState.services)
	if err != nil {
		return nil, fmt.Errorf("error initializing API: %w", err)
	}

	apiState.services = initialized

	return apiState, nil
}
