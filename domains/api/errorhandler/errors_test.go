package errorhandler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/outcatcher/anwil/domains/core/logging"
	services "github.com/outcatcher/anwil/domains/core/services/schema"
	th "github.com/outcatcher/anwil/domains/core/testhelpers"
	"github.com/stretchr/testify/require"
)

var errForTest = errors.New("magic error text")

func TestConvertErrors(t *testing.T) {
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

			recorder := th.ClosingRecorder(t)

			logWriter := bytes.Buffer{}
			logger := log.New(&logWriter, "", 0)
			ctx := logging.CtxWithLogger(context.Background(), logger)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/err/example", nil)
			require.NoError(t, err)

			echoCtx := echo.New().NewContext(req, recorder)

			HandleErrors()(data.inputErr, echoCtx)
			require.NoError(t, err)

			require.Contains(t, logWriter.String(), data.inputErr.Error())
			require.EqualValues(t, data.expectedBody, recorder.Body.String())
		})
	}
}
