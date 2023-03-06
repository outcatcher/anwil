package middlewares

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	services "github.com/outcatcher/anwil/domains/services/dto"
)

// statusCodeFromError returns status code for corresponding error.
func statusCodeFromError(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, services.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, services.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, services.ErrNotFound), errors.Is(err, sql.ErrNoRows):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// ConvertErrors converts response error to valid status code.
func ConvertErrors(c *gin.Context) {
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

	c.AbortWithStatus(statusCode)

	log.Printf("error response code %d (%s) with reason: %s", statusCode, http.StatusText(statusCode), err)
}
