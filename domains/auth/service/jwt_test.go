package service

import (
	"encoding/hex"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/outcatcher/anwil/domains/auth/dto"
	"github.com/outcatcher/anwil/domains/auth/service/schema"
	services "github.com/outcatcher/anwil/domains/internals/services/schema"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	privateKey = "e3de69265ea200c17633b8b7ba90c17c15e96f3f1d0ad608d9f628e515c7e53b" +
		"d6507afe638ea0565709842d869581edfc5e5b6186a8215f6bed2504991ff9fb"

	// token with claims:
	// &Claims{Username: "random-username"}.
	token = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9." +
		"eyJ1c2VybmFtZSI6InJhbmRvbS11c2VybmFtZSJ9." +
		"nVxpOGoHA9cggY3yGY9RJdZLYdPYnkBTClfG5HTLtLA4uEUKz5tdjlKvGHr0DkT9AVq1tiaC1SxC1ICcV4wECg"
)

type AuthTests struct {
	suite.Suite

	auth schema.Service
}

func (s *AuthTests) TestGenerateToken() {
	t := s.T()
	t.Parallel()

	t.Run("w/ claims", func(t *testing.T) {
		t.Parallel()

		claims := &dto.Claims{Username: "random-username"}

		tok, err := s.auth.GenerateToken(claims)
		require.NoError(t, err)
		require.NotEmpty(t, tok)
	})

	t.Run("w/o claims", func(t *testing.T) {
		t.Parallel()

		tok, err := s.auth.GenerateToken(nil)
		require.NoError(t, err)
		require.NotEmpty(t, tok)
	})

	t.Run("w/ invalid key", func(t *testing.T) {
		t.Parallel()

		auth2 := &auth{
			privateKey: make([]byte, 0),
		}

		_, err := auth2.GenerateToken(nil)
		require.ErrorIs(t, err, dto.ErrInvalidPrivateKeySize)
	})
}

func (s *AuthTests) TestValidateToken() {
	t := s.T()

	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		// claims := &Claims{Username: "random-username"}

		claims, err := s.auth.ValidateToken(token)
		require.NoError(t, err)
		require.Equal(t, "random-username", claims.Username)
	})

	t.Run("invalid key", func(t *testing.T) {
		t.Parallel()

		privateKey, err := hex.DecodeString("27a2fd4868ca3c71dbecfb8f89c75f48de642d95f4efbfe47bf401b7" +
			"8935b0786ec428ede4c0d6cba5d12fe166c67b660177f879a4bb750ee67dceec1b624eee")
		require.NoError(t, err)

		auth2 := &auth{
			privateKey: privateKey,
		}

		_, err = auth2.ValidateToken(token)
		require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
		require.ErrorIs(t, err, services.ErrUnauthorized)
	})

	t.Run("invalid algorithm", func(t *testing.T) {
		t.Parallel()

		tok := jwt.New(jwt.SigningMethodHS512)

		privateKeyDecoded, err := hex.DecodeString(privateKey)
		require.NoError(t, err)

		signedString, err := tok.SignedString(privateKeyDecoded)
		require.NoError(t, err)

		_, err = s.auth.ValidateToken(signedString)
		require.ErrorIs(t, err, dto.ErrUnexpectedSignMethod)
		require.ErrorIs(t, err, services.ErrUnauthorized)
	})
}

func TestAuth(t *testing.T) {
	t.Parallel()

	privateKeyDecoded, err := hex.DecodeString(privateKey)
	require.NoError(t, err)

	s := &AuthTests{
		auth: &auth{
			privateKey: privateKeyDecoded,
		},
	}

	t.Log(hex.EncodeToString(privateKeyDecoded))

	suite.Run(t, s)
}
