/*
Package services provides general means of service management and initialization
*/
package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/outcatcher/anwil/domains/core/services/schema"
)

var (
	errCyclicServiceDependency = errors.New("service dependency cycle detected")
	errDefinitionMissing       = errors.New("service definition is missing")

	errNotNeeded   = errors.New("consumer not expecting injection")
	errNotProvided = errors.New("provider doesn't provide injection")
)

type serviceState int

const (
	_ serviceState = iota // not started
	serviceInProgress
	serviceReady
)

// initializer is a helper storing initialization process state.
type initializer struct {
	state any

	serviceDefinitions map[schema.ServiceID]schema.ServiceDefinition
	services           map[schema.ServiceID]any
	serviceStates      map[schema.ServiceID]serviceState
}

func (init *initializer) initWithDependencies(
	ctx context.Context,
	id schema.ServiceID,
) error {
	svc, ok := init.serviceDefinitions[id]
	if !ok {
		return fmt.Errorf("service %s: %w", id, errDefinitionMissing) // impossible in current implementation
	}

	switch init.serviceStates[id] {
	case serviceInProgress:
		return fmt.Errorf("service %s: %w", id, errCyclicServiceDependency)
	case serviceReady:
		return nil
	}

	init.serviceStates[id] = serviceInProgress

	for _, depID := range svc.DependsOn {
		if err := init.initWithDependencies(ctx, depID); err != nil {
			return err
		}
	}

	initialized, err := svc.Init(ctx, init.state)
	if err != nil {
		return fmt.Errorf("error initializing service %s: %w", id, err)
	}

	init.services[id] = initialized
	init.serviceStates[id] = serviceReady

	return nil
}

// Initialize initializes given services with given state.
//
// Service dependencies will be checked for existing cycles and initialized in the dependency order.
//
// State will be passed to each service in mapping `Init` method.
func Initialize(
	ctx context.Context, state any, services ...schema.ServiceDefinition,
) (schema.ServiceMapping, error) {
	serviceDefMap := make(map[schema.ServiceID]schema.ServiceDefinition, len(services))

	for _, def := range services {
		serviceDefMap[def.ID] = def
	}

	initer := initializer{
		state:              state,
		serviceDefinitions: serviceDefMap,
		services:           make(map[schema.ServiceID]any, len(services)),
		serviceStates:      make(map[schema.ServiceID]serviceState),
	}

	for _, service := range services {
		if err := initer.initWithDependencies(ctx, service.ID); err != nil {
			return nil, err
		}
	}

	return initer.services, nil
}

// InjectFunc - function injecting something into service.
type InjectFunc func(consumer, provider any) error

// ValidateArgInterfaces is a helper method to assert types of consumer and provider.
//
// Example
//
//	reqStorage, provStorage, err := services.ValidateArgInterfaces[RequiresStorage, WithStorage](serv, state)
func ValidateArgInterfaces[TCons any, TProv any](consumer, provider any) (TCons, TProv, error) {
	var (
		cons TCons
		prov TProv
		ok   bool
	)

	cons, ok = consumer.(TCons)
	if !ok {
		return cons, prov, errNotNeeded
	}

	prov, ok = provider.(TProv)
	if !ok {
		return cons, prov, errNotProvided
	}

	return cons, prov, nil
}

// InjectServiceWith - initialize service with given service inject functions.
//
// Inject functions expected to add some other service reference into given service.
//
// Example:
//
//	err := services.InjectServiceWith(
//		authService, state,
//		configDTO.ConfigInject,
//		logSchema.LoggerInject,
//	)
func InjectServiceWith(service, state any, injects ...InjectFunc) error {
	for _, initFunc := range injects {
		if err := initFunc(service, state); err != nil {
			return fmt.Errorf("error intializing service: %w", err)
		}
	}

	return nil
}
