package testing

import (
	"net/http"

	"github.com/stretchr/testify/require"
)

func (s *AnwilSuite) TestEcho() {
	t := s.T()

	response := s.request("GET", "/api/v1/echo", nil, nil)

	require.Equal(t, http.StatusOK, response.Code)
	require.EqualValues(t, response.Body.Bytes(), []byte("OK"))
}
