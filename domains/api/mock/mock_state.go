/*
Package mock contains mock implementation of state.
*/
package mock

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/outcatcher/anwil/domains/auth"
	authDTO "github.com/outcatcher/anwil/domains/auth/dto"
	configDTO "github.com/outcatcher/anwil/domains/config/dto"
	"github.com/outcatcher/anwil/domains/storage"
	storageDTO "github.com/outcatcher/anwil/domains/storage/dto"
	"github.com/outcatcher/anwil/domains/users"
	usersDTO "github.com/outcatcher/anwil/domains/users/dto"
	"github.com/stretchr/testify/require"
)

var mockTMPDir = filepath.Join(os.TempDir(), "anwil_mock")

type mockState struct {
	cfg *configDTO.Configuration

	t *testing.T

	storage *storage.MockDBExecutor

	logger *log.Logger

	auth  authDTO.Service
	users usersDTO.Service
}

// Config returns API configuration.
func (m *mockState) Config() *configDTO.Configuration {
	return m.cfg
}

// Authentication returns auth service instance.
func (m *mockState) Authentication() authDTO.Service {
	return m.auth
}

// Users returns users service instance.
func (m *mockState) Users() usersDTO.Service {
	return m.users
}

// Storage returns storage.
func (m *mockState) Storage() storageDTO.QueryExecutor {
	return m.storage
}

// Logger returns logger instance.
func (m *mockState) Logger() *log.Logger {
	return m.logger
}

func getRandomPort(t *testing.T) int {
	t.Helper()

	nBig, err := rand.Int(rand.Reader, big.NewInt(0xffff)) //nolint:gomnd
	require.NoError(t, err)

	return int(nBig.Int64())
}

func generatePrivateKeyFile(t *testing.T) string {
	t.Helper()

	_, key, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	tmpFile, err := os.CreateTemp(mockTMPDir, "*_key")
	require.NoError(t, err)

	defer func() {
		require.NoError(t, tmpFile.Close())
	}()

	t.Cleanup(func() {
		require.NoError(t, os.Remove(tmpFile.Name()))
	})

	keyString := hex.EncodeToString(key)

	_, err = io.WriteString(tmpFile, keyString)
	require.NoError(t, err)

	return tmpFile.Name()
}

// initServices is simplified services initialization.
func (m *mockState) initServices(ctx context.Context) {
	m.auth = auth.New()
	m.users = users.New()

	require.NoError(m.t, m.auth.Init(ctx, m))
	require.NoError(m.t, m.users.Init(ctx, m))
}

// NewAPIMock generates new API state for testing purposes.
func NewAPIMock(t *testing.T) *mockState {
	t.Helper()

	err := os.MkdirAll(mockTMPDir, os.ModePerm)
	require.NoError(t, err)

	cfg := &configDTO.Configuration{
		API: configDTO.APIConfiguration{
			Host:       "mock-api-host",
			Port:       getRandomPort(t),
			StaticPath: "mock-static-path",
		},
		DB: configDTO.DatabaseConfiguration{
			Host:          "mock-db-host",
			Port:          getRandomPort(t),
			DatabaseName:  "mock-db",
			Username:      "mock-username",
			Password:      "mock-password",
			MigrationsDir: "mock-migrations",
		},
		PrivateKeyPath: generatePrivateKeyFile(t),
		Debug:          true,
	}

	state := &mockState{cfg: cfg, t: t, logger: log.Default()}

	state.initServices(context.Background())

	return state
}
