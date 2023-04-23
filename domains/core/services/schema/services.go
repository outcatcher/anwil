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
	// Unique service ID allowing to build service dependency tree
	ID ServiceID
	// Init is a function used to initialized service instance with given state
	Init ServiceInitFunc
	// DependsOn list IDs of the required services.
	DependsOn []ServiceID
	// InitHandlersFunc is a function for initializing service-related handlers
	// after the service is initialized.
	// Can remain `nil` if service does not provide API endpoints.
	InitHandlersFunc func(state ProvidingServices) AddHandlersFunc
}

// ServiceMapping - ServiceID to Service mapping.
type ServiceMapping map[ServiceID]any

// ProvidingServices describes provider of the services.
type ProvidingServices interface {
	Service(id ServiceID) any
}
