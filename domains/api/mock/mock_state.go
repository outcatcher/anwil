package mock

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/domains/auth"
	authDTO "github.com/outcatcher/anwil/domains/auth/dto"
	configDTO "github.com/outcatcher/anwil/domains/config/dto"
	"github.com/outcatcher/anwil/domains/storage"
	storageDTO "github.com/outcatcher/anwil/domains/storage/dto"
	"github.com/outcatcher/anwil/domains/users"
	usersDTO "github.com/outcatcher/anwil/domains/users/dto"
	"github.com/stretchr/testify/require"
)

var (
	mockTMPDir = filepath.Join(os.TempDir(), "anwil_mock")

	errNotImplemented = errors.New("mock missing method implementation")
)

type mockState struct {
	cfg *configDTO.Configuration

	t *testing.T

	storage *storage.MockDBExecutor

	auth  authDTO.Service
	users usersDTO.Service
}

func (m *mockState) Config() *configDTO.Configuration {
	return m.cfg
}

func (m *mockState) Authentication() authDTO.Service {
	return m.auth
}

func (m *mockState) Users() usersDTO.Service {
	return m.users
}

func (m *mockState) Storage() storageDTO.QueryExecutor {
	return m.storage
}

func (m *mockState) Serve(context.Context) (*http.Server, error) {
	return nil, fmt.Errorf("(Serve): %w", errNotImplemented)
}

func (m *mockState) NewRouter(context.Context, ...gin.HandlerFunc) (*gin.Engine, error) {
	return nil, fmt.Errorf("(NewRouter): %w", errNotImplemented)
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
func (m *mockState) initServices() {
	m.auth = auth.New()
	m.users = users.New()

	require.NoError(m.t, m.auth.Init(m))
	require.NoError(m.t, m.users.Init(m))
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

	state := &mockState{cfg: cfg, t: t}

	state.initServices()

	return state
}
