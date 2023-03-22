//go:build acceptance

package testing

import (
	"net/http"

	"github.com/stretchr/testify/require"
)

func (s *AnwilSuite) TestLogin() {
	t := s.T()
	t.Parallel()

	s.login()
}

func (s *AnwilSuite) TestLoginInvalidPassword() {
	t := s.T()
	t.Parallel()

	response := s.requestJSON(
		http.MethodPost,
		parseRequestURL(t, "/api/v1/login"),
		map[string]interface{}{
			"username": debugUsername,
			"password": "asdafqfqwef!",
		},
		nil,
	)

	require.Equal(t, http.StatusUnauthorized, response.Code)
}

func (s *AnwilSuite) TestLoginMissingCredentials() {
	t := s.T()
	t.Parallel()

	response := s.requestJSON(
		http.MethodPost,
		parseRequestURL(t, "/api/v1/login"),
		mapBody{},
		nil,
	)

	require.Equal(t, http.StatusBadRequest, response.Code)
}
