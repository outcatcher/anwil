package middlewares

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	services "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/stretchr/testify/require"
)

const defaultTimeoutMs = 100

var errForTest = errors.New("magic error text")

type withLogger struct {
	l *log.Logger
}

func (w *withLogger) Logger() *log.Logger {
	return w.l
}

func newLoggerState(output io.Writer) *withLogger {
	return &withLogger{log.New(output, "", log.LstdFlags)}
}

func TestConvertErrors(t *testing.T) { //nolint:funlen
	t.Parallel()

	cases := []struct {
		inputErr     error
		expectedCode int
		expectedBody string
	}{
		{
			services.ErrUnauthorized,
			http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized),
		},
		{
			services.ErrForbidden,
			http.StatusForbidden,
			http.StatusText(http.StatusForbidden),
		},
		{
			errForTest,
			http.StatusInternalServerError,
			errForTest.Error(),
		},
		{
			services.ErrConflict,
			http.StatusConflict,
			services.ErrConflict.Error(),
		},
		{
			services.ErrNotFound,
			http.StatusNotFound,
			services.ErrNotFound.Error(),
		},
	}

	for _, data := range cases {
		data := data

		t.Run(fmt.Sprint(data.expectedCode), func(t *testing.T) {
			t.Parallel()

			logWriter := new(bytes.Buffer)
			state := newLoggerState(logWriter)

			app := fiber.New(fiber.Config{ErrorHandler: ConvertErrors(state)})

			app.Get("/", func(c *fiber.Ctx) error {
				return data.inputErr
			})

			resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil), defaultTimeoutMs)
			require.NoError(t, err)

			require.EqualValues(t, data.expectedCode, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.EqualValues(t, data.expectedBody, string(respBody))

			require.Contains(t, logWriter.String(), data.inputErr.Error())
		})
	}
}
