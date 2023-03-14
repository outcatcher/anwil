package dto

import (
	"errors"
	"fmt"

	services "github.com/outcatcher/anwil/domains/services/dto"
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

// InitWithAuth initializes given service with authentication.
func InitWithAuth(service interface{}, state interface{}) error {
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
	services.Service

	EncryptPassword(src string) (string, error)
	ValidatePassword(input, encrypted string) error
	ValidateToken(token string) (*Claims, error)
	GenerateToken(claims *Claims) (string, error)
}
