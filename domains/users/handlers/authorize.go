package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	users "github.com/outcatcher/anwil/domains/users/dto"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

type credentialsRequest struct {
	// Username
	Username string `json:"username" validate:"required"`
	// User password
	Password string `json:"password" validate:"required"`
}

type jwtResponse struct {
	Token string `json:"token"`
}

func handleAuthorize(usr schema.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(credentialsRequest)

		if err := bindAndValidateJSON(c, req); err != nil {
			return fmt.Errorf("error authorizing user: %w", err)
		}

		user := users.User{
			Username: req.Username,
			Password: req.Password,
		}

		tok, err := usr.GenerateUserToken(c.Context(), user)
		if err != nil {
			return fmt.Errorf("error authorizing user: %w", err)
		}

		_ = c.SendStatus(http.StatusOK)

		return c.JSON(jwtResponse{Token: tok})
	}
}
