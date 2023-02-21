// We're using closingRecorder which closes body during cleanup
//
//nolint:bodyclose
package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/api/ctxhelpers"
	"github.com/outcatcher/anwil/internal/auth"
	"github.com/stretchr/testify/require"
)

func TestJWTMiddlewareAuthorized(t *testing.T) {
	t.Parallel()

	recorder := &httptest.ResponseRecorder{}
	ginCtx, _ := gin.CreateTestContext(recorder)

	username := "random-user-name!"

	privateKey, err := auth.LoadPrivateKey("") // will generate private key
	require.NoError(t, err)

	tokenString, err := auth.GenerateToken(&auth.Claims{Username: username}, privateKey)
	require.NoError(t, err)

	ginCtx.Request = &http.Request{
		Header: map[string][]string{
			authHeader: {fmt.Sprintf("Bearer %s", tokenString)},
		},
	}

	JWTAuth(privateKey)(ginCtx)

	require.Len(t, ginCtx.Errors, 0)
	require.Equal(t, http.StatusOK, recorder.Result().StatusCode)
	require.EqualValues(t, username, ginCtx.Value(ctxhelpers.CtxKeyUsername))
}

func TestJWTMiddlewareMissingHeader(t *testing.T) {
	t.Parallel()

	recorder := &httptest.ResponseRecorder{}
	ginCtx, _ := gin.CreateTestContext(recorder)

	expectedKey, err := auth.LoadPrivateKey("") // will generate private key
	require.NoError(t, err)

	ginCtx.Request = &http.Request{
		Header: make(map[string][]string),
	}

	JWTAuth(expectedKey)(ginCtx)

	require.Len(t, ginCtx.Errors, 0)
	require.Equal(t, http.StatusUnauthorized, recorder.Result().StatusCode)
	require.Equal(t, nil, ginCtx.Value(ctxhelpers.CtxKeyUsername))
}

func closingRecorder(t *testing.T) *httptest.ResponseRecorder {
	t.Helper()

	recorder := new(httptest.ResponseRecorder)

	t.Cleanup(func() {
		result := recorder.Result()
		if result == nil {
			return
		}

		_ = result.Body.Close()
	})

	return recorder
}

func TestJWTMiddlewareInvalidSign(t *testing.T) {
	t.Parallel()

	recorder := closingRecorder(t)
	ginCtx, _ := gin.CreateTestContext(recorder)

	signKey, err := auth.LoadPrivateKey("") // will generate private key
	require.NoError(t, err)

	tokenString, err := auth.GenerateToken(&auth.Claims{}, signKey)
	require.NoError(t, err)

	ginCtx.Request = &http.Request{
		Header: map[string][]string{
			authHeader: {fmt.Sprintf("Bearer %s", tokenString)},
		},
	}

	// server and sign key differs
	serverKey, err := auth.LoadPrivateKey("") // will generate private key
	require.NoError(t, err)

	JWTAuth(serverKey)(ginCtx)

	require.Len(t, ginCtx.Errors, 0)
	require.Equal(t, http.StatusForbidden, recorder.Result().StatusCode)
	require.Equal(t, nil, ginCtx.Value(ctxhelpers.CtxKeyUsername))
}
