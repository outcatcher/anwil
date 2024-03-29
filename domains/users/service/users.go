package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/outcatcher/anwil/domains/core/errbase"
	"github.com/outcatcher/anwil/domains/users/service/schema"
	"github.com/outcatcher/anwil/domains/users/storage"
)

// GetUser returns user data by username.
func (u *service) GetUser(ctx context.Context, username string) (*schema.User, error) {
	user, err := u.storage.GetUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &schema.User{
		UUID:     user.UUID,
		Username: user.Username,
		Password: user.Password,
		FullName: user.FullName,
	}, nil
}

// SaveUser saves new user data.
//
// user.Password expected to be not encrypted.
func (u *service) SaveUser(ctx context.Context, user schema.User) error {
	if err := checkRequirements(user.Password); err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	pwd, err := encrypt(user.Password, u.privateKey)
	if err != nil {
		return fmt.Errorf("error encrypting new user password: %w", err)
	}

	_, err = u.GetUser(ctx, user.Username)

	switch {
	case errors.Is(err, errbase.ErrNotFound):
	case err == nil:
		return fmt.Errorf("%w: user %s already exist", errbase.ErrConflict, user.Username)
	default:
		return err
	}

	err = u.storage.InsertUser(ctx, storage.Wisher{
		Username: user.Username,
		Password: pwd,
		FullName: user.FullName,
	})
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}

// GenerateUserToken validates user credentials and returns token.
func (u *service) GenerateUserToken(ctx context.Context, user schema.User) (string, error) {
	existing, err := u.GetUser(ctx, user.Username)
	if errors.Is(err, errbase.ErrNotFound) {
		return "", fmt.Errorf("user %s: %w", user.Username, errbase.ErrNotFound)
	}

	if err != nil {
		return "", fmt.Errorf("error retreating user: %w", err)
	}

	err = validatePassword(user.Password, existing.Password, u.privateKey)
	if err != nil {
		return "", fmt.Errorf("error validating user credentials: %w", err)
	}

	jwtClaims := &schema.Claims{
		Username: user.Username,
		UserUUID: user.Username,
	}

	tok, err := Generate(jwtClaims, u.privateKey)
	if err != nil {
		return "", fmt.Errorf("error generating user token: %w", err)
	}

	return tok, nil
}
