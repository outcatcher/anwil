/*
Package dto contains DTOs describing services and
*/
package dto

import "fmt"

// Service - base interface for service to adhere.
type Service interface {
	Init(state interface{}) error
}

type initializer func(service, state interface{}) error

// InitializeWith - initialize service with given initializers.
func InitializeWith(service Service, state interface{}, inits ...initializer) error {
	for _, initFunc := range inits {
		if err := initFunc(service, state); err != nil {
			return fmt.Errorf("error intializing service: %w", err)
		}
	}

	return nil
}
