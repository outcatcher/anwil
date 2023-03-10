package dto

import (
	"errors"
)

var (
	// ErrNotFound - error for requested entity being missing.
	ErrNotFound = errors.New("not found")
	// ErrForbidden - error for requested operation to be forbidden.
	ErrForbidden = errors.New("forbidden")
	// ErrUnauthorized - error for authorization failures.
	ErrUnauthorized = errors.New("not authorized")
)
