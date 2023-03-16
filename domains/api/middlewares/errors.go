package middlewares

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/internals/logging"
	services "github.com/outcatcher/anwil/domains/internals/services/schema"
)

type responseError struct {
	Code   int    `json:"-"`
	Reason string `json:"reason,omitempty"`
}

// statusCodeFromError returns status code for corresponding error.
func statusCodeFromError(err error) responseError {
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

	err := c.Errors[0]
	if err == nil {
		return
	}

	statusCode := statusCodeFromError(err)

	if statusCode.Reason == "" {
		c.AbortWithStatus(statusCode.Code)
	} else {
		c.AbortWithStatusJSON(statusCode.Code, statusCode)
	}

	log.Printf("error response code %d (%s) with reason: %s",
		statusCode.Code, http.StatusText(statusCode.Code), err,
	)
}
