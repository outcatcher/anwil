package users

import (
	"context"
	"fmt"
	"log"

	authDTO "github.com/outcatcher/anwil/domains/auth/service/schema"
	"github.com/outcatcher/anwil/domains/internals/services"
	svcDTO "github.com/outcatcher/anwil/domains/internals/services/schema"
	storageDTO "github.com/outcatcher/anwil/domains/internals/storage/schema"
	"github.com/outcatcher/anwil/domains/users/schema"
	userStorage "github.com/outcatcher/anwil/domains/users/storage"
)

type users struct {
	storage *userStorage.UserStorage

	log *log.Logger

	auth authDTO.Service
}

// UseAuthentication - use given service as an auth service for users.
func (u *users) UseAuthentication(auth authDTO.Service) {
	u.auth = auth
}

// UseStorage attaches given DB storage to the service.
func (u *users) UseStorage(db storageDTO.QueryExecutor) {
	u.storage = userStorage.New(db)
}

// UseLogger attaches logger to the service.
func (u *users) UseLogger(logger *log.Logger) {
	u.log = logger
}

// DependsOn defines services Users service depends on.
func (*users) DependsOn() []svcDTO.ServiceID {
	return []svcDTO.ServiceID{authDTO.ServiceAuth}
}

// ID returns  users service ID.
func (*users) ID() svcDTO.ServiceID {
	return schema.ServiceUsers
}

// Init initialized user instance with given state.
func (u *users) Init(_ context.Context, state interface{}) error {
	err := services.InjectServiceWith(
		u, state,
		storageDTO.StorageInject,
		authDTO.AuthInject,
	)
	if err != nil {
		return fmt.Errorf("erorr initializing user service: %w", err)
	}

	return nil
}

// New returns not initialized instance of users service.
func New() schema.Service {
	return new(users)
}
