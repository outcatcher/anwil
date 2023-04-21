package errorhandler

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/outcatcher/anwil/domains/core/errbase"
	th "github.com/outcatcher/anwil/domains/core/testhelpers"
	"github.com/stretchr/testify/require"
)

var (
	errForTest = errors.New("magic error text")
	logWriter  = &bytes.Buffer{}
)

func TestMain(m *testing.M) {
	log.SetOutput(logWriter)

	os.Exit(m.Run())
}

func TestConvertErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		inputErr     error
		expectedCode int
		expectedBody string
	}{
		{
			errbase.ErrUnauthorized,
			http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized),
		},
		{
			errbase.ErrForbidden,
			http.StatusForbidden,
			http.StatusText(http.StatusForbidden),
		},
		{
			errForTest,
			http.StatusInternalServerError,
			errForTest.Error(),
		},
		{
			errbase.ErrConflict,
			http.StatusConflict,
			errbase.ErrConflict.Error(),
		},
		{
			errbase.ErrNotFound,
			http.StatusNotFound,
			errbase.ErrNotFound.Error(),
		},
	}

	for _, data := range cases {
		data := data

		t.Run(fmt.Sprint(data.expectedCode), func(t *testing.T) {
			t.Parallel()

			recorder := th.ClosingRecorder(t)

			url := fmt.Sprintf("/err/example/%s", th.RandomString("", 5))
			method := http.MethodGet

			req, err := http.NewRequest(method, url, nil)
			require.NoError(t, err)

			echoCtx := echo.New().NewContext(req, recorder)

			HandleErrors()(data.inputErr, echoCtx)
			require.NoError(t, err)

			expectedErrString := fmt.Sprintf("Error performing %s %s: %s", method, url, data.inputErr.Error())

			require.Contains(t, logWriter.String(), expectedErrString)
			require.EqualValues(t, data.expectedBody, recorder.Body.String())
		})
	}
}
