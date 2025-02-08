package apierror

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

func (e *APIError) ToResponse() (int, string) {
	return e.Code, e.Message
}
