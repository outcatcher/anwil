package middlewares

import (
	"net/http"
	"net/http/httptest"
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

	okStatus := http.StatusIMUsed

	for _, method := range cases {
		method := method

		t.Run(method, func(t *testing.T) {
			t.Parallel()

			app := fiber.New(fiber.Config{ErrorHandler: func(c *fiber.Ctx, err error) error {
				return c.SendStatus(http.StatusInternalServerError)
			}})
			app.Get("/", RequireJSON, func(c *fiber.Ctx) error {
				return c.SendStatus(okStatus)
			})

			resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil), 100)
			require.NoError(t, err)

			require.EqualValues(t, okStatus, resp.StatusCode)
		})
	}
}

func TestRequireJSONOk(t *testing.T) {
	t.Parallel()

	okStatus := http.StatusIMUsed

	cases := []string{http.MethodPut, http.MethodPost}

	for _, method := range cases {
		method := method

		t.Run(method, func(t *testing.T) {
			t.Parallel()

			app := fiber.New(fiber.Config{ErrorHandler: func(c *fiber.Ctx, err error) error {
				return c.SendStatus(http.StatusInternalServerError)
			}})
			app.Get("/", RequireJSON, func(c *fiber.Ctx) error {
				return c.SendStatus(okStatus)
			})

			resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil), 100)
			require.NoError(t, err)

			require.EqualValues(t, okStatus, resp.StatusCode)
		})
	}
}
