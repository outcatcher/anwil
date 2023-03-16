package password

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptPassword(t *testing.T) {
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
