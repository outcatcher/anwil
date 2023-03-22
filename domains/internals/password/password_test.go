package password

import (
	"crypto/rand"
	"testing"

	services "github.com/outcatcher/anwil/domains/internals/services/schema"
	"github.com/stretchr/testify/require"
)

func TestPasswordWorkflow(t *testing.T) {
	t.Parallel()

	randomBytes := make([]byte, 128)

	_, err := rand.Read(randomBytes)
	require.NoError(t, err)

	inputPassword := "truly-random-password"

	encrypted, err := Encrypt(inputPassword, randomBytes)
	require.NoError(t, err)

	require.NotEmpty(t, encrypted)
	require.Len(t, encrypted, 128) // sha512 encrypted string is 128 bytes long

	require.NoError(t, Validate(inputPassword, encrypted, randomBytes))
}

func TestValidate_invalid(t *testing.T) {
	t.Parallel()

	randomBytes := make([]byte, 128)

	_, err := rand.Read(randomBytes)
	require.NoError(t, err)

	inputPassword := "truly-random-password"

	encrypted, err := Encrypt(inputPassword, randomBytes)
	require.NoError(t, err)

	require.NotEmpty(t, encrypted)
	require.Len(t, encrypted, 128) // sha512 encrypted string is 128 bytes long

	err = Validate(inputPassword+"no!", encrypted, randomBytes)
	require.ErrorIs(t, err, services.ErrUnauthorized)
}

func TestEncode_noSecret(t *testing.T) {
	t.Parallel()

	_, err := Encrypt("", nil)
	require.ErrorIs(t, err, errMissingPrivateKey)
}

func TestValidate_noSecret(t *testing.T) {
	t.Parallel()

	err := Validate("", "", nil)
	require.ErrorIs(t, err, errMissingPrivateKey)
}
