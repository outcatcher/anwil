package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func handleEcho(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).SendString("OK")
}
