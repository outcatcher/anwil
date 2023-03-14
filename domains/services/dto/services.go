/*
Package dto contains DTOs describing services and
*/
package dto

import (
	"context"
	"fmt"
)

// ServiceID - ID of the service.
type ServiceID string

// Service - base interface for service to adhere.
type Service interface {
	// ID returns unique service ID.
	ID() ServiceID

	// Init initialized service instance with given state.
	Init(ctx context.Context, state interface{}) error

	// DependsOn lists services the service depends on
	DependsOn() []ServiceID
}

// ServiceMapping - ID to service mapping.
type ServiceMapping map[ServiceID]Service

type serviceInit func(service, state interface{}) error

// InitializeWith - initialize service with given service initializers.
func InitializeWith(service Service, state interface{}, inits ...serviceInit) error {
	for _, initFunc := range inits {
		if err := initFunc(service, state); err != nil {
			return fmt.Errorf("error intializing service: %w", err)
		}
	}

	return nil
}
