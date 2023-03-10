package auth

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"log"

	"github.com/outcatcher/anwil/domains/auth/dto"
	configDTO "github.com/outcatcher/anwil/domains/config/dto"
	logDTO "github.com/outcatcher/anwil/domains/logging/dto"
	services "github.com/outcatcher/anwil/domains/services/dto"
)

type auth struct {
	cfg *configDTO.Configuration

	log *log.Logger

	privateKey ed25519.PrivateKey
}

// UseConfig attaches configuration to the service.
func (a *auth) UseConfig(configuration *configDTO.Configuration) {
	a.cfg = configuration
}

// UseLogger attaches logger to the service.
func (a *auth) UseLogger(logger *log.Logger) {
	a.log = logger
}

// Init initializes new auth service.
func (a *auth) Init(ctx context.Context, state interface{}) error {
	err := services.InitializeWith(
		a, state,
		configDTO.InitWithConfig,
		logDTO.InitWithLogger,
	)
	if err != nil {
		return fmt.Errorf("error initializing auth service: %w", err)
	}

	key, err := a.cfg.GetPrivateKey(ctx)
	if err != nil {
		return fmt.Errorf("error creating Auth service: %w", err)
	}

	a.privateKey = key

	return nil
}

// New returns not initialized instance of auth service.
func New() dto.Service {
	return new(auth)
}
