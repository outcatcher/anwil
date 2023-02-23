package dto

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

// Errors produced by jwt package.
var (
	ErrUnexpectedSignMethod  = errors.New("unexpected signing method")
	ErrInvalidPrivateKeySize = errors.New("private key size is invalid")
)

// Claims - JWT payload contents.
type Claims struct {
	jwt.RegisteredClaims

	Username string `json:"username"`
}
