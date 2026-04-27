package jwtutil

import "errors"

var (
	ErrInvalidUsername    = errors.New("invalid username in token")
	ErrInvalidPermissions = errors.New("invalid permissions in token")
)
