package middlewares

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/labstack/echo/v4"
	th "github.com/outcatcher/anwil/domains/core/testhelpers"
	"github.com/outcatcher/anwil/domains/core/validation"
	"github.com/stretchr/testify/require"
)

func okResponse(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func TestRequireJSONMissingHeader(t *testing.T) {
	t.Parallel()

	cases := []string{http.MethodPut, http.MethodPost}

	for _, method := range cases {
		method := method

		t.Run(method, func(t *testing.T) {
			t.Parallel()

			rec := th.ClosingRecorder(t)
			req := &http.Request{
				URL:    new(url.URL),
				Method: method,
				Header: make(http.Header),
			}

			echoCtx := echo.New().NewContext(req, rec)

			err := RequireJSON(okResponse)(echoCtx)
			require.Error(t, err)

			require.ErrorIs(t, err, validation.ErrValidationFailed)
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

			recorder := th.ClosingRecorder(t)
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

			recorder := th.ClosingRecorder(t)

			header := make(http.Header)
			header.Set("content-type", echo.MIMEApplicationJSON)

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
