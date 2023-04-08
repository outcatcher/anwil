package middlewares

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	services "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/outcatcher/anwil/domains/core/validation"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

const (
	tokenPrefix = "Bearer "

	contextKeyUsername = "username"
)

var (
	errBindingHeaders  = errors.New("error binding request headers")
	errValidateHeaders = errors.New("error validating request header")

	headerBinder = new(echo.DefaultBinder)
)

type reqAuth struct {
	Authorization string `header:"Authorization" validate:"required,jwt-header"`
}

// JWTAuth check JWT and loads user info into Gin context.
//
// This middleware happens before request is processed, so we need to abort context early,
// so main handler won't be triggered.
func JWTAuth(state schema.WithUsers) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := new(reqAuth)

			if err := bindAndValidateHeader(c, req); err != nil {
				return fmt.Errorf("%w: %w", services.ErrUnauthorized, err)
			}

			tokenString := strings.TrimPrefix(req.Authorization, tokenPrefix)

			claims, err := state.Users().ValidateUserToken(c.Request().Context(), tokenString)
			if err != nil {
				return fmt.Errorf("%w: %w", services.ErrUnauthorized, err)
			}

			c.Set(contextKeyUsername, claims.Username)

			return next(c)
		}
	}
}

// bindAndValidateHeader binds request body to structure,
// validates it using `validate` tag and trows errors into gin context.
func bindAndValidateHeader(c echo.Context, req any) error {
	authHeader := c.Request().Header.Get("Authorization")

	log.Println(authHeader)

	if err := headerBinder.BindHeaders(c, req); err != nil {
		return fmt.Errorf("%w to %T", errBindingHeaders, req)
	}

	if err := validation.ValidateHeaderCtx(c.Request().Context(), req); err != nil {
		return fmt.Errorf("%w %+v", errValidateHeaders, req)
	}

	return nil
}
