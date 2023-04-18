package middlewares

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	configSchema "github.com/outcatcher/anwil/domains/core/config/schema"
	th "github.com/outcatcher/anwil/domains/core/testhelpers"
	"github.com/outcatcher/anwil/domains/users/service"
	"github.com/outcatcher/anwil/domains/users/service/schema"
	"github.com/stretchr/testify/require"
)

type configState struct {
	cfg *configSchema.Configuration

	pKey ed25519.PrivateKey
}

func (c configState) Config() *configSchema.Configuration {
	return c.cfg
}

func newStateWithTMPKey(t *testing.T) *configState {
	t.Helper()

	_, key, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	created, err := os.CreateTemp("", "anwil-cfg-*")
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, os.Remove(created.Name()))
	})

	_, err = hex.NewEncoder(created).Write(key)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, created.Close())
	}()

	cfg := &configSchema.Configuration{
		PrivateKeyPath: created.Name(),
	}

	return &configState{cfg: cfg, pKey: key}
}

func TestJWTAuth_401(t *testing.T) {
	t.Parallel()

	state := newStateWithTMPKey(t)

	rec := th.ClosingRecorder(t)
	req := &http.Request{
		URL:    new(url.URL),
		Method: http.MethodGet,
		Header: make(http.Header),
	}

	echoCtx := echo.New().NewContext(req, rec)

	err := JWTAuth(state)(func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})(echoCtx)
	require.ErrorAs(t, err, &echo.ErrUnauthorized)
}

func TestJWTAuth_200(t *testing.T) {
	t.Parallel()

	state := newStateWithTMPKey(t)

	rec := th.ClosingRecorder(t)
	req := &http.Request{
		URL:    new(url.URL),
		Method: http.MethodGet,
		Header: make(http.Header),
	}

	username := th.RandomString("user-", 5)

	tok, err := service.Generate(&schema.Claims{Username: username}, state.pKey)
	require.NoError(t, err)

	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("%s %s", "Bearer", tok))

	echoCtx := echo.New().NewContext(req, rec)

	err = JWTAuth(state)(func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})(echoCtx)
	require.NoError(t, err)
}
