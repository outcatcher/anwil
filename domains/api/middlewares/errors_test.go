package middlewares

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/outcatcher/anwil/domains/core/logging"
	services "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/stretchr/testify/require"
)

var errForTest = errors.New("magic error text")

func TestConvertErrors(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.ReleaseMode)

	cases := []struct {
		inputErr     error
		expectedCode int
		expectedBody string
	}{
		{services.ErrUnauthorized, http.StatusUnauthorized, ""},
		{services.ErrForbidden, http.StatusForbidden, ""},

		{
			errForTest,
			http.StatusInternalServerError,
			errForTest.Error(),
		},
		{
			services.ErrConflict,
			http.StatusInternalServerError,
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

			recorder := closingRecorder(t)

			logWriter := bytes.Buffer{}
			logger := log.New(&logWriter, "", 0)
			ctx := logging.CtxWithLogger(context.Background(), logger)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/err/example", nil)
			require.NoError(t, err)

			echoCtx := echo.New().NewContext(req, recorder)

			err = ConvertErrors(func(_ echo.Context) error {
				return data.inputErr
			})(echoCtx)
			require.NoError(t, err)

			require.Contains(t, logWriter.String(), data.inputErr.Error())
			require.EqualValues(t, data.expectedBody, recorder.Body.String())
		})
	}
}
