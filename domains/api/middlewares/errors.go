package middlewares

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/outcatcher/anwil/domains/core/logging"
	services "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/outcatcher/anwil/domains/core/validation"
)

// errToHTTPError returns status code for corresponding error.
func errToHTTPError(err error) *echo.HTTPError {
	bindErr := new(echo.BindingError)
	httpError := new(echo.HTTPError)

	switch {
	case err == nil:
		return nil
	case errors.As(err, &httpError):
		return httpError
	case errors.Is(err, services.ErrUnauthorized):
		return &echo.HTTPError{Code: http.StatusUnauthorized}
	case errors.Is(err, services.ErrForbidden):
		return &echo.HTTPError{Code: http.StatusForbidden}
	case errors.Is(err, services.ErrNotFound), errors.Is(err, sql.ErrNoRows):
		return &echo.HTTPError{Code: http.StatusNotFound, Message: err.Error()}
	case errors.Is(err, services.ErrConflict):
		return &echo.HTTPError{Code: http.StatusConflict, Message: err.Error()}
	case errors.Is(err, validation.ErrValidationFailed),
		errors.As(err, &bindErr):
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	default:
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}
}

// ConvertErrors converts response error to valid status code.
//
// TODO: move to c.HTTPErrorHandler instead middleware.
func ConvertErrors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logging.LoggerFromCtx(c.Request().Context())

		// after handling
		err := next(c)
		if err == nil {
			return nil
		}

		httpError := errToHTTPError(err)

		log.Printf("Error performing %s %s: %s", c.Request().Method, c.Request().URL, err.Error())

		if httpError.Message == nil {
			return c.NoContent(httpError.Code)
		}

		return c.String(httpError.Code, fmt.Sprint(httpError.Message))
	}
}
