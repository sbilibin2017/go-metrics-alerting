package apiclient

import "fmt"

var (
	ErrTimeout         = fmt.Errorf("timeout error")
	ErrRequestFailed   = fmt.Errorf("request failed")
	ErrInvalidResponse = fmt.Errorf("invalid response")
)
