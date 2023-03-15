/*
Package services provides general means of service management and initialization
*/
package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/outcatcher/anwil/domains/services/dto"
)

var (
	errCyclicServiceDependency = errors.New("service dependency cycle detected")
	errServiceNotInState       = errors.New("service is missing in service map")
)

type serviceState int

const (
	_ serviceState = iota
	serviceInProgress
	serviceReady
)

type initializer struct {
	state interface{}

	services      dto.ServiceMapping
	serviceStates map[dto.ServiceID]serviceState
}

func (init *initializer) initWithDependencies(
	ctx context.Context,
	id dto.ServiceID,
) error {
	svc, ok := init.services[id]
	if !ok {
		return fmt.Errorf("service %s: %w", id, errServiceNotInState) // impossible in current implementation
	}

	svcState := init.serviceStates[id]

	if svcState == serviceInProgress {
		return fmt.Errorf("service %s: %w", id, errCyclicServiceDependency)
	}

	if svcState == serviceReady {
		return nil // already initialized
	}

	init.serviceStates[id] = serviceInProgress

	dependencies := svc.DependsOn()

	for _, depID := range dependencies {
		if err := init.initWithDependencies(ctx, depID); err != nil {
			return err
		}
	}

	if err := svc.Init(ctx, init.state); err != nil {
		return fmt.Errorf("error initializing service %s: %w", id, err)
	}

	init.serviceStates[id] = serviceReady

	return nil
}

// Initialize initializes given services with given state.
//
// Service dependencies will be checked for existing cycles and initialized in the dependency order.
//
// State will be passed to each service in mapping `Init` method.
func Initialize(ctx context.Context, state interface{}, svcMapping dto.ServiceMapping) (dto.ServiceMapping, error) {
	initer := initializer{
		state:         state,
		services:      svcMapping,
		serviceStates: make(map[dto.ServiceID]serviceState),
	}

	for id := range svcMapping {
		if err := initer.initWithDependencies(ctx, id); err != nil {
			return nil, err
		}
	}

	return initer.services, nil
}

// InjectFunc - function injecting something into service.
type InjectFunc func(service, state interface{}) error

// InjectServiceWith - initialize service with given service inject functions.
//
// Inject functions expected to add some other service reference into given service.
//
// Example:
//
//	err := services.InjectServiceWith(
//		authService, state,
//		configDTO.ConfigInject,
//		logDTO.LoggerInject,
//	)
func InjectServiceWith(service dto.Service, state interface{}, injects ...InjectFunc) error {
	for _, initFunc := range injects {
		if err := initFunc(service, state); err != nil {
			return fmt.Errorf("error intializing service: %w", err)
		}
	}

	return nil
}
