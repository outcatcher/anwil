package dto

import (
	"errors"
)

var ErrInvalidState = errors.New("state not matching requirements")
