//go:build acceptance

package testing

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/stretchr/testify/require"
)

func (s *AnwilSuite) TestEcho() {
	t := s.T()
	t.Parallel()

	response := s.request("GET", parseRequestURL(t, "/api/v1/echo"), nil, nil)

	require.Equal(t, http.StatusOK, response.Code)
	require.EqualValues(t, response.Body.Bytes(), []byte("OK"))
}

func (s *AnwilSuite) login() string {
	t := s.T()

	t.Helper()

	response := s.requestJSON(
		http.MethodPost,
		parseRequestURL(t, "/api/v1/login"),
		map[string]interface{}{
			"username": debugUsername,
			"password": debugPassword,
		},
		nil,
	)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.Code, string(body))

	var loginResponse struct {
		Token string `json:"token"`
	}

	require.NoError(t, json.Unmarshal(body, &loginResponse))

	return loginResponse.Token
}

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

func addAuthHeader(token string, src map[string]string) map[string]string {
	if src == nil {
		src = make(map[string]string)
	}

	src["Authorization"] = fmt.Sprintf("Bearer %s", token)
	return src
}

func (s *AnwilSuite) TestSecureEcho() {
	t := s.T()
	t.Parallel()

	token := s.login()

	response := s.requestJSON(
		http.MethodGet,
		parseRequestURL(t, "/api/v1/auth-echo"),
		nil,
		addAuthHeader(token, nil),
	)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.Code, string(body))
}
