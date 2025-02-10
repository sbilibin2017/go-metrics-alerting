package errors

import (
	e "errors"
)

// Ошибки валидации.
var (
	ErrInvalidKeyFormat = e.New("invalid key format")
	ErrValueNotFound    = e.New("value not found")
	ErrContextDone      = e.New("context canceled or timed out")
	ErrValueNotSaved    = e.New("value is not saved")
)
