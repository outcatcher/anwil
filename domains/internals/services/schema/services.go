/*
Package schema contains DTOs describing services and
*/
package schema

import (
	"context"
)

// ServiceID - ID of the service.
type ServiceID string

// Service - base interface for service.
type Service interface {
	// ID returns unique service ID.
	ID() ServiceID

	// Init initialized service instance with given state.
	//
	// Init is a place where all service dependencies are injected.
	// It should be called after all required services are loaded into state.
	Init(ctx context.Context, state interface{}) error

	// DependsOn list IDs of required services.
	DependsOn() []ServiceID
}

// ServiceMapping - ServiceID to Service mapping.
type ServiceMapping map[ServiceID]Service
