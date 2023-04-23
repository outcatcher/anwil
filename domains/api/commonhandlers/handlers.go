/*
Package commonhandlers contains unscoped API handlers
*/
package commonhandlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func handleEcho(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

// AddEchoHandlers adds common echoing endpoints.
func AddEchoHandlers(baseGroup, secGroup *echo.Group) error {
	baseGroup.GET("/echo", handleEcho)
	secGroup.GET("/auth-echo", handleEcho)

	return nil
}
