/*
Package schema contains DTOs for logging helpers
*/
package schema

import (
	"fmt"
	"log"

	"github.com/outcatcher/anwil/domains/internals/services"
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
func LoggerInject(consumer, provider any) error {
	reqLogger, provLogger, err := services.ValidateArgInterfaces[RequiresLogger, WithLogger](consumer, provider)
	if err != nil {
		return fmt.Errorf("error injecting logger: %w", err)
	}

	reqLogger.UseLogger(provLogger.Logger())

	return nil
}
