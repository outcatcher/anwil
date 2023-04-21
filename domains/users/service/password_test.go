package service

import (
	"crypto/rand"
	"testing"

	"github.com/outcatcher/anwil/domains/core/errbase"
	"github.com/stretchr/testify/require"
)

func TestPasswordWorkflow(t *testing.T) {
	t.Parallel()

	randomBytes := make([]byte, 128)

	_, err := rand.Read(randomBytes)
	require.NoError(t, err)

	inputPassword := "truly-random-password-new"

	encrypted, err := encrypt(inputPassword, randomBytes)
	require.NoError(t, err)

	require.NotEmpty(t, encrypted)
	require.Len(t, encrypted, 128) // sha512 encrypted string is 128 bytes long

	require.NoError(t, validatePassword(inputPassword, encrypted, randomBytes))
}

func TestValidate_invalid(t *testing.T) {
	t.Parallel()

	randomBytes := make([]byte, 128)

	_, err := rand.Read(randomBytes)
	require.NoError(t, err)

	inputPassword := "truly-random-password"

	encrypted, err := encrypt(inputPassword, randomBytes)
	require.NoError(t, err)

	require.NotEmpty(t, encrypted)
	require.Len(t, encrypted, 128) // sha512 encrypted string is 128 bytes long

	err = validatePassword(inputPassword+"no!", encrypted, randomBytes)
	require.ErrorIs(t, err, errbase.ErrUnauthorized)
}

func TestValidate_invalidEncrypted(t *testing.T) {
	t.Parallel()

	randomBytes := make([]byte, 128)

	_, err := rand.Read(randomBytes)
	require.NoError(t, err)

	inputPassword := "truly-random-password"

	err = validatePassword(inputPassword, inputPassword, randomBytes)
	require.Error(t, err)
}

func TestEncode_noSecret(t *testing.T) {
	t.Parallel()

	_, err := encrypt("", nil)
	require.ErrorIs(t, err, errMissingPrivateKey)
}

func TestValidate_noSecret(t *testing.T) {
	t.Parallel()

	err := validatePassword("", "", nil)
	require.ErrorIs(t, err, errMissingPrivateKey)
}
