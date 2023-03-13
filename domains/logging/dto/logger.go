/*
Package dto contains DTOs for logging helpers
*/
package dto

import (
	"errors"
	"fmt"
	"log"
)

var (
	errStateWithoutLogger   = errors.New("given state has no logging")
	errServiceWithoutLogger = errors.New("given service does not support logging")
)

// WithLogger containing logger.
type WithLogger interface {
	Logger() *log.Logger
}

// RequiresLogger can use logger.
type RequiresLogger interface {
	UseLogger(logger *log.Logger)
}

// InitWithLogger attaches logger to given service.
func InitWithLogger(service interface{}, state interface{}) error {
	reqConfig, ok := service.(RequiresLogger)
	if !ok {
		return fmt.Errorf("error intializing logging: %w", errServiceWithoutLogger)
	}

	stateWithConfig, ok := state.(WithLogger)
	if !ok {
		return fmt.Errorf("error intializing logging: %w", errStateWithoutLogger)
	}

	reqConfig.UseLogger(stateWithConfig.Logger())

	return nil
}
