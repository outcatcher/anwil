//go:build acceptance

package testing

import (
	"io"
	"net/http"
	"testing"

	th "github.com/outcatcher/anwil/domains/core/testhelpers"
	"github.com/stretchr/testify/require"
)

func randomUserData() mapBody {
	return mapBody{
		"username":  th.RandomString("user-", 5),
		"password":  th.RandomString("pwd-", 5),
		"full_name": th.RandomString("I AM ", 8),
	}
}

func (s *AnwilSuite) TestUserCreate_200() {
	t := s.T()

	t.Parallel()

	resp := s.requestJSON(
		http.MethodPost,
		parseRequestURL(t, "/api/v1/wisher"),
		randomUserData(),
		nil,
	)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.EqualValues(t, http.StatusCreated, resp.StatusCode, body)
}

func (s *AnwilSuite) TestUserCreate_400() {
	t := s.T()

	t.Parallel()

	cases := map[string]mapBody{
		"missing_fields": {},
		"invalid_type": {
			"username":  123,
			"password":  16,
			"full_name": 2032,
		},
		"too_simple_password": {
			"username": th.RandomString("user", 5),
			"password": "qwertyui123!",
		},
	}

	for name, data := range cases {
		data := data

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := s.requestJSON(
				http.MethodPost,
				parseRequestURL(t, "/api/v1/wisher"),
				data,
				nil,
			)

			require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func (s *AnwilSuite) TestUserCreate_409() {
	t := s.T()

	t.Parallel()

	// debug user already exists
	userData := mapBody{
		"username":  debugUsername,
		"password":  debugPassword,
		"full_name": debugFullName,
	}

	resp := s.requestJSON(
		http.MethodPost,
		parseRequestURL(t, "/api/v1/wisher"),
		userData,
		nil,
	)
	require.EqualValues(t, http.StatusConflict, resp.StatusCode)
}
