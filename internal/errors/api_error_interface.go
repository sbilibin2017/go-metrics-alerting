package errors

type ApiErrorInterface interface {
	Status() int   // Код состояния HTTP, связанный с ошибкой.
	Error() string // Описание ошибки.
}
