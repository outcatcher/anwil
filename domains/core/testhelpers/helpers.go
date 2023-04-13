/*
Package testhelpers contains various test helper functions
*/
package testhelpers

import (
	"bytes"
	"crypto/rand"
	"net/http/httptest"
	"testing"
)

const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// RandomString generates random string with given prefix.
func RandomString(prefix string, length int) string {
	bytes := make([]byte, length)

	_, _ = rand.Read(bytes)

	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}

	return prefix + string(bytes)
}

// ClosingRecorder creates new httptest.ResponseRecorder and closing body at cleanup.
func ClosingRecorder(t *testing.T) *httptest.ResponseRecorder {
	t.Helper()

	recorder := &httptest.ResponseRecorder{Body: new(bytes.Buffer)}

	t.Cleanup(func() {
		result := recorder.Result()
		if result == nil {
			return
		}

		_ = result.Body.Close()
	})

	return recorder
}
