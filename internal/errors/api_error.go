package errors

type ApiErrorInterface interface {
	Status() int
	Error() string // Implementing the Error method from the `error` interface
}

type ApiError struct {
	StatusCode int
	Message    string
}

func (e *ApiError) Status() int {
	return e.StatusCode
}

// Implementing the Error method from the `error` interface
func (e *ApiError) Error() string {
	return e.Message
}

var _ ApiErrorInterface = &ApiError{}
