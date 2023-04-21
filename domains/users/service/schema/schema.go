/*
Package schema contains service definition for Users service
*/
package schema

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/outcatcher/anwil/domains/core/services/schema"
)

// ServiceID - ID for user service.
const ServiceID schema.ServiceID = "users"

// Claims - JWT payload contents.
type Claims struct {
	jwt.RegisteredClaims

	Username string `json:"username"`
	UserUUID string `json:"user_uuid"`
}

var (
	// ErrUnexpectedSignMethod - signing method not supported.
	ErrUnexpectedSignMethod = errors.New("unexpected signing method")
	// ErrInvalidPrivateKeySize - ed25519 private key size not matched.
	ErrInvalidPrivateKeySize = errors.New("private key size is invalid")
)

// UserService - service handling user-related functionality.
type UserService interface {
	GetUser(ctx context.Context, username string) (*User, error)
	SaveUser(ctx context.Context, user User) error
	GenerateUserToken(ctx context.Context, user User) (string, error)
}

// User holds user data.
type User struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"-"` // hex-encoded password, make sure it's not reaching JSON
	FullName string `json:"full_name"`
}
