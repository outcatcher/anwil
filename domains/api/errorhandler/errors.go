/*
Package errorhandler contains handler for handler-produced errors
*/
package errorhandler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/outcatcher/anwil/domains/core/errbase"
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
	case errors.Is(err, errbase.ErrUnauthorized):
		return echo.ErrUnauthorized
	case errors.Is(err, errbase.ErrForbidden):
		return echo.ErrForbidden
	case errors.Is(err, errbase.ErrNotFound), errors.Is(err, sql.ErrNoRows):
		return &echo.HTTPError{Code: http.StatusNotFound, Message: err.Error()}
	case errors.Is(err, errbase.ErrConflict):
		return &echo.HTTPError{Code: http.StatusConflict, Message: err.Error()}
	case errors.Is(err, validation.ErrValidationFailed),
		errors.As(err, &bindErr):
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	default:
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}
}

// HandleErrors converts response error to valid status code.
func HandleErrors() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		httpError := errToHTTPError(err)

		log.Printf("Error performing %s %s: %s", c.Request().Method, c.Request().URL, err.Error())

		responseErr := c.String(httpError.Code, fmt.Sprint(httpError.Message))
		if responseErr != nil {
			log.Printf("Error handling error %s: %s", err.Error(), responseErr.Error())
		}
	}
}
