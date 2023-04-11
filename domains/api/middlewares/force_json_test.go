package middlewares

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/outcatcher/anwil/domains/core/validation"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestRequireJSONMissingHeader(t *testing.T) {
	t.Parallel()

	cases := []string{http.MethodPut, http.MethodPost}

	for _, method := range cases {
		method := method

		t.Run(method, func(t *testing.T) {
			t.Parallel()

			app := fiber.New()

			fiberCtx := app.AcquireCtx(&fasthttp.RequestCtx{})
			fiberCtx.Method(method)
			fiberCtx.Body()

			err := RequireJSON(fiberCtx)
			require.ErrorIs(t, err, errInvalidMIMEType)
			require.ErrorIs(t, err, validation.ErrValidationFailed)
		})
	}
}

func TestRequireJSONNoContentTypeOk(t *testing.T) {
	t.Parallel()

	cases := []string{http.MethodHead, http.MethodGet, http.MethodDelete}

	for _, method := range cases {
		method := method

		t.Run(method, func(t *testing.T) {
			t.Parallel()

			app := fiber.New()

			fiberCtx := app.AcquireCtx(&fasthttp.RequestCtx{})
			fiberCtx.Method(method)
			fiberCtx.Body()

			err := RequireJSON(fiberCtx)
			require.NoError(t, err)
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

			app := fiber.New()

			fiberCtx := app.AcquireCtx(&fasthttp.RequestCtx{})
			fiberCtx.Method(method)
			fiberCtx.Request().Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			err := RequireJSON(fiberCtx)
			require.NoError(t, err)
		})
	}
}
