/*
Package schema contains DTOs for logging helpers
*/
package schema

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

// LoggerInject injects logger into service.
func LoggerInject(service interface{}, state interface{}) error {
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
