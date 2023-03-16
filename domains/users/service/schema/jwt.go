package schema

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

// Claims - JWT payload contents.
type Claims struct {
	jwt.RegisteredClaims

	Username string `json:"username"`
}

var (
	// ErrUnexpectedSignMethod - signing method not supported.
	ErrUnexpectedSignMethod = errors.New("unexpected signing method")
	// ErrInvalidPrivateKeySize - ed25519 private key size not matched.
	ErrInvalidPrivateKeySize = errors.New("private key size is invalid")
)