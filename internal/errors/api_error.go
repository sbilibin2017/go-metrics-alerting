package errors

// ApiError представляет собой структуру ошибки API, содержащую код состояния и сообщение об ошибке.
type ApiError struct {
	StatusCode int    // Код состояния HTTP, связанный с ошибкой.
	Message    string // Описание ошибки.
}

// Status возвращает HTTP-код состояния, связанный с ошибкой.
func (e *ApiError) Status() int {
	return e.StatusCode
}

// Error реализует метод интерфейса `error` и возвращает текстовое сообщение об ошибке.
func (e *ApiError) Error() string {
	return e.Message
}

// Проверка, что ApiError реализует интерфейс ApiErrorInterface.
var _ ApiErrorInterface = &ApiError{}
