/*
Package schema contains schema for Auth service.
*/
package schema

import (
	"errors"
	"fmt"

	"github.com/outcatcher/anwil/domains/auth/dto"
	svcDTO "github.com/outcatcher/anwil/domains/internals/services/schema"
)

// ServiceAuth - ID for auth service.
const ServiceAuth = "auth"

var (
	errStateWithoutAuth   = errors.New("given state has no auth")
	errServiceWithoutAuth = errors.New("given service does not support auth")
)

// WithAuth can return auth service instance.
type WithAuth interface {
	Authentication() Service
}

// requiresAuth can return auth service instance.
type requiresAuth interface {
	// UseAuthentication - use given service as auth service.
	UseAuthentication(auth Service)
}

// AuthInject injects authentication into given service.
func AuthInject(service interface{}, state interface{}) error {
	reqAuth, ok := service.(requiresAuth)
	if !ok {
		return fmt.Errorf("error intializing service auth: %w", errStateWithoutAuth)
	}

	stateWithAuth, ok := state.(WithAuth)
	if !ok {
		return fmt.Errorf("error intializing service auth: %w", errServiceWithoutAuth)
	}

	reqAuth.UseAuthentication(stateWithAuth.Authentication())

	return nil
}

// Service contains auth service definition.
type Service interface {
	svcDTO.Service

	EncryptPassword(src string) (string, error)
	ValidatePassword(input, encrypted string) error
	ValidateToken(token string) (*dto.Claims, error)
	GenerateToken(claims *dto.Claims) (string, error)
}
