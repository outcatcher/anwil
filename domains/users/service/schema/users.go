/*
Package schema contains service definition for Users service
*/
package schema

import (
	"context"

	"github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/outcatcher/anwil/domains/users/dto"
)

// ServiceUsers - ID for user service.
const ServiceUsers schema.ServiceID = "users"

// Service is definition of user service.
type Service interface {
	schema.Service

	GetUser(ctx context.Context, username string) (*dto.User, error)
	SaveUser(ctx context.Context, user dto.User) error
	GenerateUserToken(ctx context.Context, user dto.User) (string, error)
	ValidateUserToken(ctx context.Context, token string) (*dto.JWTClaims, error)
}

// WithUsers can return users service instance.
type WithUsers interface {
	Users() Service
}
