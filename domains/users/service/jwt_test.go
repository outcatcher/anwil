package service

import (
	"crypto"
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/outcatcher/anwil/domains/core/errbase"
	"github.com/outcatcher/anwil/domains/users/service/schema"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	// token with claims:
	// &Claims{Username: "random-username", UserUUID: "8738ec06-7aa8-44b3-90d4-baaaf261c968"}.
	token = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9." +
		"eyJpc3MiOiJhbndpbCIsImV4cCI6MTY4NDQ5NjY1MSwiaWF0IjoxNjgxOTA0NjUxLCJ1c2VybmFtZSI6InJhbmRvbS11c2" +
		"VybmFtZSIsInVzZXJfdXVpZCI6Ijg3MzhlYzA2LTdhYTgtNDRiMy05MGQ0LWJhYWFmMjYxYzk2OCJ9." +
		"EWf57iXIFMwEr2gsFmjDiL2bYai8_1PXv9mIap411twj2-F4VPIFxv3PAeyWWKfRw-qU4RGCahuQ9CM3bD9CBg"
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

		claims := &schema.Claims{Username: "random-username", UserUUID: "8738ec06-7aa8-44b3-90d4-baaaf261c968"}

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

		claims, err := validateToken(token, s.publicKey)
		require.NoError(t, err)
		require.Equal(t, "random-username", claims.Username)
		require.Equal(t, "8738ec06-7aa8-44b3-90d4-baaaf261c968", claims.UserUUID)
	})

	t.Run("invalid key", func(t *testing.T) {
		t.Parallel()

		privateKey, err := hex.DecodeString("27a2fd4868ca3c71dbecfb8f89c75f48de642d95f4efbfe47bf401b7" +
			"8935b0786ec428ede4c0d6cba5d12fe166c67b660177f879a4bb750ee67dceec1b624eee")
		require.NoError(t, err)

		_, err = validateToken(token, ed25519.PrivateKey(privateKey).Public())
		require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
		require.ErrorIs(t, err, errbase.ErrUnauthorized)
	})

	t.Run("invalid algorithm", func(t *testing.T) {
		t.Parallel()

		tok := jwt.New(jwt.SigningMethodHS512)

		signedString, err := tok.SignedString([]byte(s.privateKey))
		require.NoError(t, err)

		_, err = validateToken(signedString, s.publicKey)
		require.ErrorIs(t, err, schema.ErrUnexpectedSignMethod)
		require.ErrorIs(t, err, errbase.ErrUnauthorized)
	})
}

func TestAuth(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(AuthTests))
}
