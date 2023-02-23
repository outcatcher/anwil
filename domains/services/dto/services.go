/*
Package dto contains DTOs describing services and
*/
package dto

import "fmt"

type Service interface {
	Init(state interface{}) error
}

type Initializer func(service, state interface{}) error

func InitializeWith(service Service, state interface{}, inits ...Initializer) error {
	for _, initFunc := range inits {
		if err := initFunc(service, state); err != nil {
			return fmt.Errorf("error intializing service: %w", err)
		}
	}

	return nil
}
