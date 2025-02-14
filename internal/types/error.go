package types

type APIErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
