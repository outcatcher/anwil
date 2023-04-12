package middlewares

import (
	"context"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/outcatcher/anwil/domains/core/config/schema"
	"github.com/outcatcher/anwil/domains/users/token"
)

const (
	contextKeyUsername = "username"
)

// JWTAuth check JWT and loads user info into Gin context.
//
// This middleware happens before request is processed, so we need to abort context early,
// so main handler won't be triggered.
func JWTAuth(ctx context.Context, state schema.WithConfig) echo.MiddlewareFunc {
	pKey, err := state.Config().GetPrivateKey(ctx)
	if err != nil {
		return nil
	}

	return func(n echo.HandlerFunc) echo.HandlerFunc {
		return echojwt.WithConfig(echojwt.Config{
			ContextKey:     contextKeyUsername,
			SigningKey:     pKey,
			KeyFunc:        token.Ed25519KeyFunc(pKey.Public()),
			ParseTokenFunc: nil,
		})(n)
	}
}
