/*
Package password contains password-related functions
*/
package password

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"

	services "github.com/outcatcher/anwil/domains/core/services/schema"
)

var errMissingPrivateKey = errors.New("missing encryption key")

func encrypt(src string, key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("error encrypting password: %w", errMissingPrivateKey)
	}

	hmc := hmac.New(sha512.New, key)

	_, err := hmc.Write([]byte(src))
	if err != nil {
		return nil, fmt.Errorf("error encrypting password: %w", err)
	}

	return hmc.Sum(nil), nil
}

// Encrypt encrypts password with a private key.
func Encrypt(src string, key []byte) (string, error) {
	encrypted, err := encrypt(src, key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encrypted), nil
}

// Validate compares given password to be equal with a given encrypted password.
func Validate(input, encrypted string, key []byte) error {
	macCompared, err := hex.DecodeString(encrypted)
	if err != nil {
		return fmt.Errorf("error decoding encrypted password: %w", err)
	}

	macInput, err := encrypt(input, key)
	if err != nil {
		return err
	}

	if !hmac.Equal(macInput, macCompared) {
		return fmt.Errorf("%w: invalid password", services.ErrUnauthorized)
	}

	return nil
}
