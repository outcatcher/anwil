package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/outcatcher/anwil/domains/core/validation"
)

// RequireJSON forces usage of application/json content type in POST and PUT request to servers.
func RequireJSON(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		contentType := c.Request().Header.Get(echo.HeaderContentType)

		if (c.Request().Method == http.MethodPost || c.Request().Method == http.MethodPut) &&
			!strings.HasPrefix(contentType, echo.MIMEApplicationJSON) {
			err := fmt.Errorf(
				"%w: invalid MIME type, %s expected",
				validation.ErrValidationFailed, echo.MIMEApplicationJSON,
			)

			return err
		}

		return next(c)
	}
}
