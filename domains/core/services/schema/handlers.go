package schema

import (
	"github.com/gofiber/fiber/v2"
)

// AddHandlersFunc - function adding handlers to secure or/and unsecure groups.
type AddHandlersFunc func(baseGroup, secGroup fiber.Router) error
