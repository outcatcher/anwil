package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	users "github.com/outcatcher/anwil/domains/users/dto"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

type reqUserCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

func handleUserRegister(usr schema.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		req := new(reqUserCreate)

		if err := c.Bind(req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)

			return
		}

		err := usr.SaveUser(ctx, users.User{
			Username: req.Username,
			Password: req.Password,
			FullName: "",
		})
		if err != nil {
			// FIXME: determine exact error to status code mapping
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}
	}
}
