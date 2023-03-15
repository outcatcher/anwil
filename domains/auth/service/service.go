package service

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"log"

	"github.com/outcatcher/anwil/domains/auth/service/schema"
	configSchema "github.com/outcatcher/anwil/domains/internals/config/schema"
	logDTO "github.com/outcatcher/anwil/domains/internals/logging/schema"
	"github.com/outcatcher/anwil/domains/internals/services"
	svcDTO "github.com/outcatcher/anwil/domains/internals/services/schema"
)

type auth struct {
	cfg *configSchema.Configuration

	log *log.Logger

	privateKey ed25519.PrivateKey
}

// ID of the auth service.
func (*auth) ID() svcDTO.ServiceID {
	return schema.ServiceAuth
}

// UseConfig attaches configuration to the service.
func (a *auth) UseConfig(configuration *configSchema.Configuration) {
	a.cfg = configuration
}

// UseLogger attaches logger to the service.
func (a *auth) UseLogger(logger *log.Logger) {
	a.log = logger
}

// DependsOn defines services auth service depends on.
func (*auth) DependsOn() []svcDTO.ServiceID {
	return []svcDTO.ServiceID{}
}

// Init initializes new auth service.
func (a *auth) Init(ctx context.Context, state interface{}) error {
	err := services.InjectServiceWith(
		a, state,
		configSchema.ConfigInject,
		logDTO.LoggerInject,
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
func New() schema.Service {
	return new(auth)
}
