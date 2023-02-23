/*
Package dto contains DTOs for User domain.
*/
package dto

import (
	"context"

	services "github.com/outcatcher/anwil/domains/services/dto"
)

// User holds user data.
type User struct {
	Username string
	Password string `json:"-"` // hex-encoded password, make sure it's not reaching JSON
	FullName string
}

// Service is definition of user service.
type Service interface {
	services.Service

	GetUser(ctx context.Context, username string) (*User, error)
	SaveUser(ctx context.Context, user User) error
}

// WithUsers can return users service instance.
type WithUsers interface {
	Users() Service
}
