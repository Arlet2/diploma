package core

import "errors"

var (
	ErrBadToken        = errors.New("incorrect jwt token")
	ErrTokenExpired    = errors.New("jwt token was expired")
	ErrAuthServerError = errors.New("auth server error")
	ErrDeviceNotFound  = errors.New("device not found in auth server")
)
