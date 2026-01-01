package server

import "errors"

var (
	ErrEmptyAuthHeader    = errors.New("empty auth header")
	ErrInternal           = errors.New("internal server error")
	ErrBadRequest         = errors.New("bad request")
	ErrInvalidTokenFormat = errors.New("invalid token format")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token expired")
)
