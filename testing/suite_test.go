package testing

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/outcatcher/anwil/api/handlers"
	"github.com/outcatcher/anwil/config"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// AnwilSuite - handlers tests.
type AnwilSuite struct {
	suite.Suite

	apiHandler http.HandlerFunc
}

func (s *AnwilSuite) request(method, target string, values url.Values, body io.Reader) *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()

	request := httptest.NewRequest(method, target, body)
	request.URL.RawQuery = values.Encode()

	s.apiHandler(responseRecorder, request)

	return responseRecorder
}

func (s *AnwilSuite) SetupSuite() {
	t := s.T()

	gin.SetMode(gin.ReleaseMode) // no need for request logs

	cfg, err := config.LoadServerConfiguration("./fixtures/test_config.yaml")
	require.NoError(t, err)

	router, err := handlers.NewRouter("..", cfg)
	require.NoError(t, err)

	// gin engine as a handler function
	s.apiHandler = func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	}
}

func TestRun(t *testing.T) {
	t.Parallel()

	testSuite := new(AnwilSuite)

	suite.Run(t, testSuite)
}
