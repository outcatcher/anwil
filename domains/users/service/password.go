package service

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"

	services "github.com/outcatcher/anwil/domains/core/services/schema"
	"github.com/outcatcher/anwil/domains/core/validation"
	pwdValidator "github.com/wagslane/go-password-validator"
)

const minEntropy = 50

var errMissingPrivateKey = errors.New("missing encryption key")

func encryptBytes(src string, key []byte) ([]byte, error) {
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

// encrypt encrypts password with a private key.
func encrypt(src string, key []byte) (string, error) {
	encrypted, err := encryptBytes(src, key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encrypted), nil
}

// validatePassword compares given password to be equal with a given encrypted password.
func validatePassword(input, encrypted string, key []byte) error {
	macCompared, err := hex.DecodeString(encrypted)
	if err != nil {
		return fmt.Errorf("error decoding encrypted password: %w", err)
	}

	macInput, err := encryptBytes(input, key)
	if err != nil {
		return err
	}

	if !hmac.Equal(macInput, macCompared) {
		return fmt.Errorf("%w: invalid password", services.ErrUnauthorized)
	}

	return nil
}

// checkRequirements validates password strength requirements.
func checkRequirements(password string) error {
	err := pwdValidator.Validate(password, minEntropy)
	if err != nil {
		return fmt.Errorf("%w: %w", validation.ErrValidationFailed, err)
	}

	return nil
}
