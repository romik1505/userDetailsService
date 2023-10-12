package service

import "errors"

var (
	ErrBadRequest    = errors.New("bad request")
	ErrEmptyResponse = errors.New("empty response")
	ErrNotFound      = errors.New("not found")
	ErrInternalError = errors.New("internal error")
)
