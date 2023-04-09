package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	users "github.com/outcatcher/anwil/domains/users/dto"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

type createUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	FullName string `json:"full_name"`
}

func handleUserRegister(usr schema.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		req := new(createUser)

		if err := bindAndValidateJSON(c, req); err != nil {
			return fmt.Errorf("error registering user: %w", err)
		}

		err := usr.SaveUser(ctx, users.User{
			Username: req.Username,
			Password: req.Password,
			FullName: req.FullName,
		})
		if err != nil {
			return fmt.Errorf("error registering user: %w", err)
		}

		return c.SendStatus(http.StatusCreated)
	}
}
