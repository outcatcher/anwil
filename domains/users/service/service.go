/*
Package service contains user service methods
*/
package service

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"log"

	configSchema "github.com/outcatcher/anwil/domains/core/config/schema"
	logSchema "github.com/outcatcher/anwil/domains/core/logging/schema"
	"github.com/outcatcher/anwil/domains/core/services"
	svcSchema "github.com/outcatcher/anwil/domains/core/services/schema"
	storageSchema "github.com/outcatcher/anwil/domains/storage/schema"
	"github.com/outcatcher/anwil/domains/users/service/schema"
	userStorage "github.com/outcatcher/anwil/domains/users/storage"
)

type users struct {
	cfg     *configSchema.Configuration
	storage userStorage.UserStorage

	log *log.Logger

	privateKey ed25519.PrivateKey
}

// UseConfig attaches configuration to the service.
func (u *users) UseConfig(configuration *configSchema.Configuration) {
	u.cfg = configuration
}

// UseStorage attaches given DB storage to the service.
func (u *users) UseStorage(db storageSchema.QueryExecutor) {
	u.storage = userStorage.New(db)
}

// UseLogger attaches logger to the service.
func (u *users) UseLogger(logger *log.Logger) {
	u.log = logger
}

// DependsOn defines services Users service depends on.
func (*users) DependsOn() []svcSchema.ServiceID {
	return []svcSchema.ServiceID{}
}

// ID returns  users service ID.
func (*users) ID() svcSchema.ServiceID {
	return schema.ServiceUsers
}

// Init initialized user instance with given state.
func (u *users) Init(ctx context.Context, state interface{}) error {
	err := services.InjectServiceWith(
		u, state,
		storageSchema.StorageInject,
		logSchema.LoggerInject,
		configSchema.ConfigInject,
	)
	if err != nil {
		return fmt.Errorf("error initializing user service: %w", err)
	}

	key, err := u.cfg.GetPrivateKey(ctx)
	if err != nil {
		return fmt.Errorf("error initializing user service: %w", err)
	}

	u.privateKey = key

	return nil
}

// New returns not initialized instance of users service.
func New() schema.Service {
	return new(users)
}
