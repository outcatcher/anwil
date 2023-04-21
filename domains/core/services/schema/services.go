/*
Package schema contains DTOs describing services and
*/
package schema

import (
	"context"
	"errors"
)

var (
	// ErrInvalidType - error casting general service to exact type.
	ErrInvalidType = errors.New("invalid service type")

	// ErrMissingService - provider has no matching service to return.
	ErrMissingService = errors.New("provider missing initialized service")
)

// ServiceID - ID of the service.
type ServiceID string

// ServiceInitFunc - function returning initialized service instance.
type ServiceInitFunc func(ctx context.Context, state any) (any, error)

// ServiceDefinition contains service metadata and initialization.
type ServiceDefinition struct {
	ID ServiceID
	// Init return initialized service instance with given state.
	Init ServiceInitFunc
	// DependsOn list IDs of required services.
	DependsOn []ServiceID
}

// ServiceMapping - ServiceID to Service mapping.
type ServiceMapping map[ServiceID]any

// ProvidingServices describes provider of the services.
type ProvidingServices interface {
	Service(id ServiceID) any
}
