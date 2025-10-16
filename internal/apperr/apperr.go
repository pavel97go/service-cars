package apperr

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal server error")
)
