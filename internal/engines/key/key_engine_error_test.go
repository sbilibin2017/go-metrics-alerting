package key

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestKeyEngineError(t *testing.T) {
	err := &KeyEngineError{msg: "test error message"}
	assert.Equal(t, "test error message", err.Error())
}

func TestErrInvalidKeyFormat(t *testing.T) {
	assert.Equal(t, "invalid key format, expected 'type:name'", ErrInvalidKeyFormat.Error())
}
