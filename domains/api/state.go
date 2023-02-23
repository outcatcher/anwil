package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/api/handlers"
	authDTO "github.com/outcatcher/anwil/domains/auth/dto"
	"github.com/outcatcher/anwil/domains/config"
	configDTO "github.com/outcatcher/anwil/domains/config/dto"
	"github.com/outcatcher/anwil/domains/storage"
	storageDTO "github.com/outcatcher/anwil/domains/storage/dto"
	usersDTO "github.com/outcatcher/anwil/domains/users/dto"
)

const defaultTimeout = time.Minute

type State struct {
	cfg *configDTO.Configuration

	storage        storageDTO.QueryExecutor
	serviceMapping map[serviceID]interface{}
}

func (s *State) Serve() (*http.Server, error) {
	cfg := s.Config()

	router, err := s.NewRouter(gin.Logger(), gin.Recovery())
	if err != nil {
		return nil, fmt.Errorf("error creating new router: %w", err)
	}

	server := &http.Server{ //nolint:exhaustruct
		Addr:              fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port),
		Handler:           router,
		ReadHeaderTimeout: defaultTimeout,
	}

	loggedAddr := server.Addr

	if cfg.API.Host == "" {
		loggedAddr = fmt.Sprintf("localhost:%d", cfg.API.Port)
	}

	log.Printf("Anwil API server started at http://%s", loggedAddr)

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

// Config returns server configuration.
func (s *State) Config() *configDTO.Configuration {
	return s.cfg
}

// Authentication service.
func (s *State) Authentication() authDTO.Service {
	return s.serviceMapping[serviceAuth].(authDTO.Service) //nolint:forcetypeassert
}

// Users service.
func (s *State) Users() usersDTO.Service {
	return s.serviceMapping[serviceUsers].(usersDTO.Service) //nolint:forcetypeassert
}

func (s *State) Storage() storageDTO.QueryExecutor {
	return s.storage
}

// Init initializes API and returns new API instance.
func Init(ctx context.Context, configPath string) (*State, error) {
	apiState := new(State)

	cfg, err := config.LoadServerConfiguration(ctx, path.Clean(configPath))
	if err != nil {
		return nil, fmt.Errorf("error loading server config: %w", err)
	}

	apiState.cfg = cfg

	db, err := storage.Connect(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the storage: %w", err)
	}

	apiState.storage = db

	if err := apiState.initServices(); err != nil {
		return nil, fmt.Errorf("error initializing API: %w", err)
	}

	return apiState, nil
}
