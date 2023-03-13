package users

import (
	"context"
	"fmt"
	"log"

	authDTO "github.com/outcatcher/anwil/domains/auth/dto"
	services "github.com/outcatcher/anwil/domains/services/dto"
	storageDTO "github.com/outcatcher/anwil/domains/storage/dto"
	"github.com/outcatcher/anwil/domains/users/dto"
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

// Init initialized user instance with given state.
func (u *users) Init(_ context.Context, state interface{}) error {
	err := services.InitializeWith(
		u, state,
		storageDTO.InitWithStorage,
		authDTO.InitWithAuth,
	)
	if err != nil {
		return fmt.Errorf("erorr initializing user service: %w", err)
	}

	return nil
}

// New returns not initialized instance of users service.
func New() dto.Service {
	return new(users)
}
