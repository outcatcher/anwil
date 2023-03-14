package middlewares

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

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

func TestRequireJSONMissingHeader(t *testing.T) {
	t.Parallel()

	cases := []string{http.MethodPut, http.MethodPost}

	for _, method := range cases {
		method := method

		t.Run(method, func(t *testing.T) {
			t.Parallel()

			recorder := closingRecorder(t)
			ginCtx, _ := gin.CreateTestContext(recorder)

			ginCtx.Request = &http.Request{
				Method: method,
				Header: make(http.Header),
			}

			RequireJSON(ginCtx)

			result := recorder.Result()
			require.NotNil(t, result)

			require.Equal(t, http.StatusBadRequest, result.StatusCode)
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

			RequireJSON(ginCtx)

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
