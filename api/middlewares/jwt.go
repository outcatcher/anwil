package middlewares

import (
	"crypto/ed25519"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/api/ctxhelpers"
	"github.com/outcatcher/anwil/internal/auth"
)

const (
	authHeader = "Authorization"

	tokenPrefix = "Bearer "
)

// JWTAuth check JWT and loads user info into Gin context.
func JWTAuth(privateKey ed25519.PrivateKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authHeader)

		if header == "" {
			log.Println("empty auth header")

			c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			c.Abort()

			return
		}

		tokenString := strings.TrimPrefix(header, tokenPrefix)

		claims, err := auth.ValidateToken(tokenString, privateKey.Public())
		if err != nil {
			log.Println("error in JWT:", err)

			c.String(http.StatusForbidden, http.StatusText(http.StatusForbidden))
			c.Abort()

			return
		}

		c.Set(ctxhelpers.CtxKeyUsername, claims.Username)
	}
}
