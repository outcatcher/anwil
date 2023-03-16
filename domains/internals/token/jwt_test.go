package token

import (
	"crypto"
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	services "github.com/outcatcher/anwil/domains/internals/services/schema"
	"github.com/outcatcher/anwil/domains/users/service/schema"
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

	privateKey ed25519.PrivateKey
	publicKey  crypto.PublicKey
}

func (s *AuthTests) SetupSuite() {
	pKey, err := hex.DecodeString(privateKey)
	require.NoError(s.T(), err)

	s.privateKey = pKey
	s.publicKey = s.privateKey.Public()
}

func (s *AuthTests) TestGenerateToken() {
	t := s.T()
	t.Parallel()

	t.Run("w/ claims", func(t *testing.T) {
		t.Parallel()

		claims := &schema.Claims{Username: "random-username"}

		tok, err := Generate(claims, s.privateKey)
		require.NoError(t, err)
		require.NotEmpty(t, tok)
	})

	t.Run("w/o claims", func(t *testing.T) {
		t.Parallel()

		tok, err := Generate(nil, s.privateKey)
		require.NoError(t, err)
		require.NotEmpty(t, tok)
	})

	t.Run("w/ invalid key", func(t *testing.T) {
		t.Parallel()

		_, err := Generate(nil, make([]byte, 0))
		require.ErrorIs(t, err, schema.ErrInvalidPrivateKeySize)
	})
}

func (s *AuthTests) TestValidateToken() {
	t := s.T()

	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		claims, err := Validate(token, s.publicKey)
		require.NoError(t, err)
		require.Equal(t, "random-username", claims.Username)
	})

	t.Run("invalid key", func(t *testing.T) {
		t.Parallel()

		privateKey, err := hex.DecodeString("27a2fd4868ca3c71dbecfb8f89c75f48de642d95f4efbfe47bf401b7" +
			"8935b0786ec428ede4c0d6cba5d12fe166c67b660177f879a4bb750ee67dceec1b624eee")
		require.NoError(t, err)

		_, err = Validate(token, ed25519.PrivateKey(privateKey).Public())
		require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
		require.ErrorIs(t, err, services.ErrUnauthorized)
	})

	t.Run("invalid algorithm", func(t *testing.T) {
		t.Parallel()

		tok := jwt.New(jwt.SigningMethodHS512)

		signedString, err := tok.SignedString([]byte(s.privateKey))
		require.NoError(t, err)

		_, err = Validate(signedString, s.publicKey)
		require.ErrorIs(t, err, schema.ErrUnexpectedSignMethod)
		require.ErrorIs(t, err, services.ErrUnauthorized)
	})
}

func TestAuth(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(AuthTests))
}
