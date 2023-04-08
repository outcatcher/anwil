package schema

import (
	"github.com/labstack/echo/v4"
)

// AddHandlersFunc - function adding handlers to secure or/and unsecure groups.
type AddHandlersFunc func(baseGroup, secGroup *echo.Group) error
