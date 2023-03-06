package auth

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/outcatcher/anwil/domains/auth/dto"
	services "github.com/outcatcher/anwil/domains/services/dto"
)

const (
	jwtDefaultExpiration = time.Hour * 24 * 30
	jwtIssuer            = "anwil"
)

// ValidateToken validates token and return JWT payload data.
func (a *auth) ValidateToken(tokenString string) (*dto.Claims, error) {
	claims := new(dto.Claims)

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodEd25519)
		if !ok {
			return nil, fmt.Errorf("%w: %s", dto.ErrUnexpectedSignMethod, token.Header["alg"])
		}

		return a.privateKey.Public(), nil
	})
	if err != nil {
		var validationErr *jwt.ValidationError

		if errors.As(err, &validationErr) {
			return nil, fmt.Errorf("error validating JWT: %w (%s)", services.ErrUnauthorized, err.Error())
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

// GenerateToken generates token with given claims.
func (a *auth) GenerateToken(claims *dto.Claims) (string, error) {
	if claims == nil {
		claims = new(dto.Claims)
	}

	if len(a.privateKey) != ed25519.PrivateKeySize {
		return "", fmt.Errorf("error generating token: %w", dto.ErrInvalidPrivateKeySize)
	}

	claims.RegisteredClaims = defaultClaims()

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	tokenString, err := token.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("error creating signed string: %w", err)
	}

	return tokenString, nil
}
