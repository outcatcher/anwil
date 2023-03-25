package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	users "github.com/outcatcher/anwil/domains/users/dto"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

type createUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	FullName string `json:"full_name"`
}

func handleUserRegister(usr schema.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		req := new(createUser)

		if err := bindAndValidateJSON(c, req); err != nil {
			return
		}

		err := usr.SaveUser(ctx, users.User{
			Username: req.Username,
			Password: req.Password,
			FullName: req.FullName,
		})
		if err != nil {
			_ = c.Error(err)

			return
		}

		c.AbortWithStatus(http.StatusCreated)
	}
}
