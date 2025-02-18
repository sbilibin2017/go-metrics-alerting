package types

type APIErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
