package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	authDTO "github.com/outcatcher/anwil/domains/auth/dto"
	services "github.com/outcatcher/anwil/domains/services/dto"
	"github.com/outcatcher/anwil/domains/users/dto"
	userStorage "github.com/outcatcher/anwil/domains/users/storage"
)

// GetUser returns user data by username.
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

// GetUserToken validates user credentials and returns token.
func (u *users) GetUserToken(ctx context.Context, user dto.User) (string, error) {
	existing, err := u.GetUser(ctx, user.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: %w", user.Username, services.ErrNotFound)
	}

	err = u.auth.ValidatePassword(user.Password, existing.Password)
	if err != nil {
		return "", fmt.Errorf("error validating user credentials: %w", err)
	}

	tok, err := u.auth.GenerateToken(&authDTO.Claims{Username: user.Username})
	if err != nil {
		return "", fmt.Errorf("error generating user token: %w", err)
	}

	return tok, nil
}
