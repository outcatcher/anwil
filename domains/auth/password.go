package auth

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	errInvalidPassword   = errors.New("invalid password")
	errMissingPrivateKey = errors.New("missing private key")
)

func encrypt(src string, key ed25519.PrivateKey) ([]byte, error) {
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

// EncryptPassword encrypts password with same key as used for JWT.
func (a *auth) EncryptPassword(src string) (string, error) {
	encrypted, err := encrypt(src, a.privateKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encrypted), nil
}

func (a *auth) ValidatePassword(input, encrypted string) error {
	macCompared, err := hex.DecodeString(encrypted)
	if err != nil {
		return fmt.Errorf("error decoding encrypted password: %w", err)
	}

	macInput, err := encrypt(input, a.privateKey)
	if err != nil {
		return err
	}

	if !hmac.Equal(macInput, macCompared) {
		return errInvalidPassword
	}

	return nil
}
