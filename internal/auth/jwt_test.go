package auth

import (
	"crypto/ed25519"
	"encoding/hex"
	"io"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
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

func checkKeyForSignature(t *testing.T, key ed25519.PrivateKey) {
	t.Helper()

	// check if JWT can use generated key for ED DSA signing
	_, err := jwt.SigningMethodEdDSA.Sign("test string", key)
	require.NoError(t, err)
}

func TestGeneratePrivateKey(t *testing.T) {
	t.Parallel()

	key, err := GeneratePrivateKey()
	require.NoError(t, err)

	checkKeyForSignature(t, key)
}

func TestLoadPrivateKey(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		key, err := LoadPrivateKey("")
		require.NoError(t, err)

		checkKeyForSignature(t, key)
	})

	t.Run("existing", func(t *testing.T) {
		t.Parallel()

		tmpFile, err := os.CreateTemp("", "key*")
		require.NoError(t, err)

		_, err = io.WriteString(tmpFile, privateKey)
		require.NoError(t, err)

		key, err := LoadPrivateKey(tmpFile.Name())
		require.NoError(t, err)

		checkKeyForSignature(t, key)
	})

	t.Run("invalid file", func(t *testing.T) {
		t.Parallel()

		tmpFile, err := os.CreateTemp("", "key*")
		require.NoError(t, err)

		_, err = io.WriteString(tmpFile, "random-invalid-data")
		require.NoError(t, err)

		_, err = LoadPrivateKey(tmpFile.Name())
		require.ErrorContains(t, err, "error decoding loaded key file")
	})

	t.Run("missing file", func(t *testing.T) {
		t.Parallel()

		_, err := LoadPrivateKey("/random-non-existing-file")
		require.ErrorContains(t, err, "error opening private key file")
	})
}

func TestGenerateToken(t *testing.T) {
	t.Parallel()

	key, err := GeneratePrivateKey()
	require.NoError(t, err)

	t.Run("w/ claims", func(t *testing.T) {
		t.Parallel()

		claims := &Claims{Username: "random-username"}

		tok, err := GenerateToken(claims, key)
		require.NoError(t, err)
		require.NotEmpty(t, tok)
	})

	t.Run("w/o claims", func(t *testing.T) {
		t.Parallel()

		tok, err := GenerateToken(nil, key)
		require.NoError(t, err)
		require.NotEmpty(t, tok)
	})

	t.Run("w/ invalid key", func(t *testing.T) {
		t.Parallel()

		_, err := GenerateToken(nil, []byte{})
		require.ErrorIs(t, err, errInvalidPrivateKeySize)
	})
}

func TestValidateToken(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		// claims := &Claims{Username: "random-username"}

		privateKeyDecoded, err := hex.DecodeString(privateKey)
		require.NoError(t, err)

		claims, err := ValidateToken(token, ed25519.PrivateKey(privateKeyDecoded).Public())
		require.NoError(t, err)
		require.Equal(t, "random-username", claims.Username)
	})

	t.Run("invalid key", func(t *testing.T) {
		t.Parallel()

		publicKeyDecoded, err := hex.DecodeString("7312ad84d80a0f57b3218820b632f06022e1fe95343b6e275d3c45b730bc2887")
		require.NoError(t, err)

		_, err = ValidateToken(token, publicKeyDecoded)
		require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
	})

	t.Run("invalid algorithm", func(t *testing.T) {
		t.Parallel()

		tok := jwt.New(jwt.SigningMethodHS512)

		privateKeyDecoded, err := hex.DecodeString(privateKey)
		require.NoError(t, err)

		signedString, err := tok.SignedString(privateKeyDecoded)
		require.NoError(t, err)

		_, err = ValidateToken(signedString, ed25519.PrivateKey(privateKeyDecoded).Public())
		require.ErrorIs(t, err, errUnexpectedSignMethod)
	})
}
