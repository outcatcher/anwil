package auth

import (
	"crypto/ed25519"
	"fmt"

	"github.com/outcatcher/anwil/domains/auth/dto"
	configDTO "github.com/outcatcher/anwil/domains/config/dto"
	services "github.com/outcatcher/anwil/domains/services/dto"
)

type auth struct {
	cfg *configDTO.Configuration

	privateKey ed25519.PrivateKey
}

// UseConfig attaches configuration to the service.
func (a *auth) UseConfig(configuration *configDTO.Configuration) {
	a.cfg = configuration
}

// Init initializes new auth service.
func (a *auth) Init(state interface{}) error {
	err := services.InitializeWith(
		a, state,
		configDTO.InitWithConfig,
	)
	if err != nil {
		return fmt.Errorf("error initializing auth service: %w", err)
	}

	key, err := a.cfg.GetPrivateKey()
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
