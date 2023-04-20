/*
Package service contains wishes service methods
*/
package service

import (
	"context"
	"fmt"
	"log"

	configSchema "github.com/outcatcher/anwil/domains/core/config/schema"
	logSchema "github.com/outcatcher/anwil/domains/core/logging/schema"
	"github.com/outcatcher/anwil/domains/core/services"
	svcSchema "github.com/outcatcher/anwil/domains/core/services/schema"
	storageSchema "github.com/outcatcher/anwil/domains/storage/schema"
	"github.com/outcatcher/anwil/domains/wishes/service/schema"
	wishStorage "github.com/outcatcher/anwil/domains/wishes/storage"
)

type service struct {
	storage *wishStorage.Storage

	log *log.Logger
}

// UseStorage attaches given DB storage to the service.
func (s *service) UseStorage(db storageSchema.QueryExecutor) {
	s.storage = wishStorage.New(db)
}

// UseLogger attaches logger to the service.
func (s *service) UseLogger(logger *log.Logger) {
	s.log = logger
}

// DependsOn defines services service depends on.
func (*service) DependsOn() []svcSchema.ServiceID {
	return []svcSchema.ServiceID{}
}

// ID returns wishes service ID.
func (*service) ID() svcSchema.ServiceID {
	return schema.ServiceWishes
}

// Init initialized service instance with given state.
func (s *service) Init(_ context.Context, state interface{}) error {
	err := services.InjectServiceWith(
		s, state,
		storageSchema.StorageInject,
		logSchema.LoggerInject,
		configSchema.ConfigInject,
	)
	if err != nil {
		return fmt.Errorf("error initializing wish service: %w", err)
	}

	return nil
}

// New creates new wishes service instance.
func New() schema.WishesService {
	return new(service)
}
