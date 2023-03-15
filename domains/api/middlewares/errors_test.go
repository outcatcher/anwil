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
	"github.com/outcatcher/anwil/domains/internals/logging"
	services "github.com/outcatcher/anwil/domains/internals/services/schema"
	"github.com/stretchr/testify/require"
)

var errForTest = errors.New("magic error text")

func TestConvertErrors(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.ReleaseMode)

	cases := []struct {
		inputErr     error
		expectedCode int
	}{
		{errForTest, http.StatusInternalServerError},
		{services.ErrUnauthorized, http.StatusUnauthorized},
		{services.ErrForbidden, http.StatusForbidden},
		{services.ErrNotFound, http.StatusNotFound},
	}

	for _, data := range cases {
		data := data

		t.Run(fmt.Sprint(data.expectedCode), func(t *testing.T) {
			t.Parallel()

			recorder := closingRecorder(t)

			logWriter := bytes.Buffer{}
			logger := log.New(&logWriter, "", 0)
			ctx := logging.CtxWithLogger(context.Background(), logger)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
			require.NoError(t, err)

			ginCtx, _ := gin.CreateTestContext(recorder)
			ginCtx.Errors = append(ginCtx.Errors, &gin.Error{Err: data.inputErr})

			ginCtx.Request = req

			ConvertErrors(ginCtx)

			require.Contains(t, logWriter.String(), data.inputErr.Error())
			t.Log(logWriter.String())

			require.True(t, ginCtx.IsAborted())
		})
	}
}
