/*
Package testhelpers contains various test helper functions
*/
package testhelpers

import (
	"crypto/rand"
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
