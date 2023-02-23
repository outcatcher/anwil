//go:build acceptance

package testing

import (
	"net/http"
)

func (s *AnwilSuite) TestUserCreate() {
	t := s.T()

	t.Parallel()

	s.requestJSON(
		http.MethodPost,
		parseRequestURL(t, "/api/v1/wisher"),
		map[string]interface{}{},
		nil,
	)
}
