package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	services "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/outcatcher/anwil/domains/core/validation"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

const (
	tokenPrefix = "Bearer "

	contextKeyUsername = "username"
)

type reqAuth struct {
	Authorization string `header:"Authorization" validate:"required,jwt-header"`
}

// JWTAuth check JWT and loads user info into Gin context.
//
// This middleware happens before request is processed, so we need to abort context early,
// so main handler won't be triggered.
func JWTAuth(state schema.WithUsers) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(reqAuth)

		if err := bindAndValidateHeader(c, req); err != nil {
			_ = c.Error(services.ErrUnauthorized)
			c.Abort()

			return
		}

		tokenString := strings.TrimPrefix(req.Authorization, tokenPrefix)

		claims, err := state.Users().ValidateUserToken(c.Request.Context(), tokenString)
		if err != nil {
			_ = c.Error(err)
			c.Abort()

			return
		}

		c.Set(contextKeyUsername, claims.Username)
	}
}

// bindAndValidateHeader binds request body to structure,
// validates it using `validate` tag and trows errors into gin context.
func bindAndValidateHeader(c *gin.Context, req any) error {
	if err := c.BindHeader(req); err != nil {
		return c.Error(err)
	}

	if err := validation.ValidateHeaderCtx(c.Request.Context(), req); err != nil {
		return c.Error(err)
	}

	return nil
}
