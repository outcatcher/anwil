//go:build testt

package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
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

func TestRequireJSONMissingHeader(t *testing.T) {
	t.Parallel()

	cases := []string{http.MethodPut, http.MethodPost}

	for _, method := range cases {
		method := method

		t.Run(method, func(t *testing.T) {
			t.Parallel()

			app := fiber.New()
			app.Use(RequireJSON)

			handler := app.Handler()

			fiberCtx := &fasthttp.RequestCtx{}
			fiberCtx.Request.SetRequestURI("/")

			handler(fiberCtx)

			require.Equal(t, http.StatusBadRequest, fiberCtx.Response.StatusCode())
		})
	}
}

func TestRequireJSONNoContentTypeOk(t *testing.T) {
	t.Parallel()

	recorder := closingRecorder(t)
	ginCtx, _ := gin.CreateTestContext(recorder)

	cases := []string{http.MethodGet, http.MethodDelete}

	for _, method := range cases {
		method := method

		t.Run(method, func(t *testing.T) {
			t.Parallel()

			ginCtx.Request = &http.Request{
				Method: method,
				Header: make(http.Header),
			}
			fiberCtx := &fasthttp.RequestCtx{}
			fiberCtx.Method()
			fiberCtx.Request.SetRequestURI("/")

			RequireJSON(fiberCtx)

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
			ginCtx, _ := gin.CreateTestContext(recorder)

			header := make(http.Header)
			header.Set("content-type", gin.MIMEJSON)

			ginCtx.Request = &http.Request{
				Method: method,
				Header: header,
			}

			RequireJSON(ginCtx)

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
