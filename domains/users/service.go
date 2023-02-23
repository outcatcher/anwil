package users

import (
	"fmt"

	authDTO "github.com/outcatcher/anwil/domains/auth/dto"
	services "github.com/outcatcher/anwil/domains/services/dto"
	storageDTO "github.com/outcatcher/anwil/domains/storage/dto"
	"github.com/outcatcher/anwil/domains/users/dto"
	userStorage "github.com/outcatcher/anwil/domains/users/storage"
)

type users struct {
	storage *userStorage.UserStorage

	auth authDTO.Service
}

func (u *users) UseAuthentication(auth authDTO.Service) {
	u.auth = auth
}

func (u *users) UseStorage(db storageDTO.QueryExecutor) {
	u.storage = userStorage.New(db)
}

func (u *users) Init(state interface{}) error {
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
