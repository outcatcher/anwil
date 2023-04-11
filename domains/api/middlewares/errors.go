package middlewares

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	logSchema "github.com/outcatcher/anwil/domains/core/logging/schema"
	services "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/outcatcher/anwil/domains/core/validation"
)

// statusCodeFromError returns status code for corresponding error.
func statusCodeFromError(err error) *fiber.Error {
	fiberErr := new(fiber.Error)

	switch {
	case err == nil:
		return nil
	case errors.As(err, &fiberErr):
		return fiberErr
	case errors.Is(err, services.ErrUnauthorized):
		return fiber.ErrUnauthorized
	case errors.Is(err, services.ErrForbidden):
		return fiber.ErrForbidden
	case errors.Is(err, services.ErrNotFound), errors.Is(err, sql.ErrNoRows):
		return &fiber.Error{Code: http.StatusNotFound, Message: err.Error()}
	case errors.Is(err, services.ErrConflict):
		return &fiber.Error{Code: http.StatusConflict, Message: err.Error()}
	case errors.Is(err, validation.ErrValidationFailed):
		return &fiber.Error{Code: http.StatusBadRequest, Message: err.Error()}
	default:
		return &fiber.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}
}

// ConvertErrors converts response error to valid status code.
//
// TODO: make it a fiber.ErrorHandler instead of middleware.
func ConvertErrors(state logSchema.WithLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := state.Logger()

		// no processing of request
		err := c.Next()
		if err == nil {
			return nil
		}

		statusCode := statusCodeFromError(err)

		log.Printf("Error performing %s %s: %s",
			c.Method(), string(c.Request().RequestURI()), err.Error(),
		)

		c.Status(statusCode.Code)

		if statusCode.Message != "" {
			return c.SendString(statusCode.Message)
		}

		return nil
	}
}
