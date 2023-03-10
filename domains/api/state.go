package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/api/handlers"
	authDTO "github.com/outcatcher/anwil/domains/auth/dto"
	"github.com/outcatcher/anwil/domains/config"
	configDTO "github.com/outcatcher/anwil/domains/config/dto"
	"github.com/outcatcher/anwil/domains/logging"
	"github.com/outcatcher/anwil/domains/storage"
	storageDTO "github.com/outcatcher/anwil/domains/storage/dto"
	usersDTO "github.com/outcatcher/anwil/domains/users/dto"
)

const defaultTimeout = time.Minute

// State holds general application state.
type State struct {
	cfg *configDTO.Configuration

	log     *log.Logger
	storage storageDTO.QueryExecutor

	serviceMapping     map[serviceID]interface{}
	serviceMappingLock sync.Mutex
}

// Server creates new API server instance.
func (s *State) Server(ctx context.Context) (*http.Server, error) {
	cfg := s.Config()

	// context is passed as BaseContext
	router, err := s.NewRouter(gin.LoggerWithWriter(s.log.Writer()), gin.Recovery()) //nolint:contextcheck
	if err != nil {
		return nil, fmt.Errorf("error creating new router: %w", err)
	}

	server := &http.Server{ //nolint:exhaustruct
		Addr:              fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port),
		Handler:           router,
		ReadHeaderTimeout: defaultTimeout,
		BaseContext:       func(_ net.Listener) context.Context { return ctx },
	}

	loggedAddr := server.Addr

	if cfg.API.Host == "" {
		loggedAddr = fmt.Sprintf("localhost:%d", cfg.API.Port)
	}

	s.log.Printf("Anwil API server started at http://%s", loggedAddr)

	return server, nil
}

// NewRouter creates new GIN engine for Anwil API.
func (s *State) NewRouter(middles ...gin.HandlerFunc) (*gin.Engine, error) {
	engine := gin.New()
	engine.Use(middles...)

	if err := handlers.PopulateEndpoints(engine, s); err != nil {
		return nil, fmt.Errorf("error populating endpoints: %w", err)
	}

	return engine, nil
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
	s.serviceMappingLock.Lock()
	defer s.serviceMappingLock.Unlock()

	return s.serviceMapping[serviceAuth].(authDTO.Service) //nolint:forcetypeassert
}

// Users service.
func (s *State) Users() usersDTO.Service {
	s.serviceMappingLock.Lock()
	defer s.serviceMappingLock.Unlock()

	return s.serviceMapping[serviceUsers].(usersDTO.Service) //nolint:forcetypeassert
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

	if err := apiState.initServices(ctx); err != nil {
		return nil, fmt.Errorf("error initializing API: %w", err)
	}

	return apiState, nil
}
