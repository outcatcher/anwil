package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func closingRecorder(t *testing.T) *httptest.ResponseRecorder {
	t.Helper()

	recorder := &httptest.ResponseRecorder{Body: new(bytes.Buffer)}

	t.Cleanup(func() {
		result := recorder.Result()
		if result == nil {
			return
		}

		_ = result.Body.Close()
	})

	return recorder
}

func okResponse(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func TestRequireJSONMissingHeader(t *testing.T) {
	t.Parallel()

	cases := []string{http.MethodPost}

	for _, method := range cases {
		method := method

		t.Run(method, func(t *testing.T) {
			t.Parallel()

			rec := closingRecorder(t)
			req := &http.Request{
				URL:    new(url.URL),
				Method: method,
				Header: make(http.Header),
			}

			echoCtx := echo.New().NewContext(req, rec)

			err := RequireJSON(okResponse)(echoCtx)
			require.Error(t, err)

			result := rec.Result()
			require.NotNil(t, result)

			require.Equal(t, http.StatusBadRequest, result.StatusCode)
		})
	}
}

func TestRequireJSONNoContentTypeOk(t *testing.T) {
	t.Parallel()

	cases := []string{http.MethodGet, http.MethodDelete}

	for _, method := range cases {
		method := method

		t.Run(method, func(t *testing.T) {
			t.Parallel()

			recorder := closingRecorder(t)
			request := &http.Request{
				Method: method,
				Header: make(http.Header),
			}

			echoCtx := echo.New().NewContext(request, recorder)

			err := RequireJSON(okResponse)(echoCtx)
			require.NoError(t, err)

			result := recorder.Result()
			require.NotNil(t, result)

			require.Equal(t, http.StatusOK, result.StatusCode)
		})
	}
}

func TestRequireJSONOk(t *testing.T) {
	t.Parallel()

	cases := []string{http.MethodPut, http.MethodPost}

	for _, method := range cases {
		method := method

		t.Run(method, func(t *testing.T) {
			t.Parallel()

			recorder := closingRecorder(t)

			header := make(http.Header)
			header.Set("content-type", gin.MIMEJSON)

			request := &http.Request{
				Method: method,
				Header: header,
			}

			echoCtx := echo.New().NewContext(request, recorder)

			err := RequireJSON(okResponse)(echoCtx)
			require.NoError(t, err)

			result := recorder.Result()
			require.NotNil(t, result)

			responseData, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			defer func() {
				require.NoError(t, result.Body.Close())
			}()

			require.Equal(t, http.StatusOK, result.StatusCode, string(responseData))
		})
	}
}
