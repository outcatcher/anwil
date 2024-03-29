package service

import (
	"crypto"
	"crypto/ed25519"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/outcatcher/anwil/domains/core/errbase"
	"github.com/outcatcher/anwil/domains/users/service/schema"
)

const (
	jwtDefaultExpiration = time.Hour * 24 * 30
	jwtIssuer            = "anwil"
)

// Ed25519KeyFunc returns a callback function to supply the key for verification.
func Ed25519KeyFunc(key crypto.PublicKey) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodEd25519)
		if !ok {
			return nil, fmt.Errorf(
				"%w: %w: %s",
				errbase.ErrUnauthorized, schema.ErrUnexpectedSignMethod, token.Header["alg"],
			)
		}

		return key, nil
	}
}

// validateToken validates token and return JWT payload data.
func validateToken(tokenString string, key crypto.PublicKey) (*schema.Claims, error) {
	claims := new(schema.Claims)

	_, err := jwt.ParseWithClaims(tokenString, claims, Ed25519KeyFunc(key))
	if err != nil {
		var validationErr *jwt.ValidationError

		if errors.As(err, &validationErr) {
			return nil, fmt.Errorf("%w: %w", errbase.ErrUnauthorized, err)
		}

		return nil, fmt.Errorf("error validating JWT: %w", err)
	}

	return claims, nil
}

func defaultClaims() jwt.RegisteredClaims {
	now := jwt.NewNumericDate(time.Now().UTC())

	return jwt.RegisteredClaims{
		Issuer:    jwtIssuer,
		ExpiresAt: jwt.NewNumericDate(now.Add(jwtDefaultExpiration)),
		IssuedAt:  now,
	}
}

// Generate generates token with given claims.
func Generate(claims *schema.Claims, key ed25519.PrivateKey) (string, error) {
	if claims == nil {
		claims = new(schema.Claims)
	}

	if len(key) != ed25519.PrivateKeySize {
		return "", fmt.Errorf("error generating token: %w", schema.ErrInvalidPrivateKeySize)
	}

	claims.RegisteredClaims = defaultClaims()

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("error creating signed string: %w", err)
	}

	return tokenString, nil
}
