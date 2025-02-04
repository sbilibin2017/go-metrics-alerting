package key

// keyEngineError представляет ошибку некорректного ключа.
type KeyEngineError struct {
	msg string
}

func (e *KeyEngineError) Error() string {
	return e.msg
}

// ErrInvalidKeyFormat — предопределенная ошибка формата ключа.
var ErrInvalidKeyFormat = &KeyEngineError{msg: "invalid key format, expected 'type:name'"}
