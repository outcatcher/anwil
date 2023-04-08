package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var errInvalidMIMEType = errors.New("invalid MIME type")

// RequireJSON forces usage of application/json content type in POST and PUT request to servers.
func RequireJSON(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		contentType := c.Request().Header.Get(echo.HeaderContentType)

		if (c.Request().Method == http.MethodPost || c.Request().Method == http.MethodPut) &&
			!strings.HasPrefix(contentType, echo.MIMEApplicationJSON) {
			err := fmt.Errorf("%w, expected %s", errInvalidMIMEType, echo.MIMEApplicationJSON)

			c.Error(err)

			return err
		}

		return next(c)
	}
}
