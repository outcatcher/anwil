/*
Package schema contains DTOs describing services and
*/
package schema

import (
	"context"
	"errors"
	"fmt"
)

// ErrInvalidType - error casting general service to exact type.
var ErrInvalidType = errors.New("invalid service type")

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

// ProvidingServices describes provider of the services.
type ProvidingServices interface {
	Service(id ServiceID) Service
}

// GetServiceFromProvider returns service of exact type by ID.
func GetServiceFromProvider[T Service](p ProvidingServices, id ServiceID) (T, error) {
	rawService := p.Service(id)

	var svc T

	svc, ok := rawService.(T)
	if !ok {
		return svc, fmt.Errorf("%w: actual service type %T", ErrInvalidType, rawService)
	}

	return svc, nil
}
