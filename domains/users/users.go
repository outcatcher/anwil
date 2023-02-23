package users

import (
	"context"
	"fmt"

	"github.com/outcatcher/anwil/domains/users/dto"
	userStorage "github.com/outcatcher/anwil/domains/users/storage"
)

func (u *users) GetUser(ctx context.Context, username string) (*dto.User, error) {
	user, err := u.storage.GetUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &dto.User{
		Username: user.Username,
		Password: user.Password,
		FullName: user.FullName,
	}, nil
}

// SaveUser saves new user data.
//
// user.Password expected to be not encrypted.
func (u *users) SaveUser(ctx context.Context, user dto.User) error {
	password, err := u.auth.EncryptPassword(user.Password)
	if err != nil {
		return fmt.Errorf("error encrypting new user password: %w", err)
	}

	err = u.storage.InsertUser(ctx, userStorage.Wisher{
		Username: user.Username,
		Password: password,
		FullName: user.FullName,
	})
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}
