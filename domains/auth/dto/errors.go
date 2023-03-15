package dto

import "errors"

var (
	// ErrUnexpectedSignMethod - signing method not supported.
	ErrUnexpectedSignMethod = errors.New("unexpected signing method")
	// ErrInvalidPrivateKeySize - ed25519 private key size not matched.
	ErrInvalidPrivateKeySize = errors.New("private key size is invalid")
)
