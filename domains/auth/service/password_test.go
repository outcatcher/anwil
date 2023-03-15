package service_test

import (
	"testing"

	"github.com/outcatcher/anwil/domains/api/mock"
	"github.com/stretchr/testify/require"
)

func TestEncryptPassword(t *testing.T) {
	t.Parallel()

	authService := mock.NewAPIMock(t).Authentication()

	encrypted, err := authService.EncryptPassword("truly-random-password")
	require.NoError(t, err)

	require.NotEmpty(t, encrypted)
	require.Len(t, encrypted, 128) // sha512 encrypted string is 128 bytes long
}
