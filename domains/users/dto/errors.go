package dto

import "errors"

var (
	ErrNoSuchUser      = errors.New("no user found by username")
	ErrInvalidPassword = errors.New("user password don't match")
)
