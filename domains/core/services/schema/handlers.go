package schema

import (
	"github.com/gofiber/fiber/v2"
)

// AddHandlersFunc - function adding handlers to secure or/and unsecure groups.
type AddHandlersFunc func(baseGroup fiber.Router) error
