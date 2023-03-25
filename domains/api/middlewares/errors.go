package middlewares

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/core/logging"
	services "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/outcatcher/anwil/domains/core/validation"
)

type responseError struct {
	Code   int    `json:"-"`
	Reason string `json:"reason,omitempty"`
}

// statusCodeFromError returns status code for corresponding error.
func statusCodeFromError(err error) responseError {
	ginErr := new(gin.Error)

	switch {
	case err == nil:
		return responseError{Code: http.StatusOK}
	case errors.Is(err, services.ErrUnauthorized):
		return responseError{Code: http.StatusUnauthorized}
	case errors.Is(err, services.ErrForbidden):
		return responseError{Code: http.StatusForbidden}
	case errors.Is(err, services.ErrNotFound), errors.Is(err, sql.ErrNoRows):
		return responseError{Code: http.StatusNotFound, Reason: err.Error()}
	case errors.Is(err, services.ErrConflict):
		return responseError{Code: http.StatusConflict, Reason: err.Error()}
	case errors.Is(err, validation.ErrValidationFailed),
		errors.As(err, &ginErr) && ginErr.IsType(gin.ErrorTypeBind):
		return responseError{Code: http.StatusBadRequest, Reason: err.Error()}
	default:
		return responseError{Code: http.StatusInternalServerError, Reason: err.Error()}
	}
}

// ConvertErrors converts response error to valid status code.
func ConvertErrors(c *gin.Context) {
	log := logging.LoggerFromCtx(c.Request.Context())

	// no processing of request
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	err := c.Errors[len(c.Errors)-1] // last error defines the code
	statusCode := statusCodeFromError(err.Err)

	errLines := make([]string, len(c.Errors))
	// log all other errors
	for i, err := range c.Errors {
		errLines[i] = err.Error()
	}

	log.Printf("Error performing %s  %s:\n\t%s",
		c.Request.Method, c.Request.URL, strings.Join(errLines, "\n\t"),
	)

	if statusCode.Reason == "" {
		c.AbortWithStatus(statusCode.Code)
	} else {
		c.AbortWithStatusJSON(statusCode.Code, statusCode)
	}
}
