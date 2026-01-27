package audit

import "fmt"

type ClientError struct {
	message string
	err     error
}

var ClientErr *ClientError

func NewClientError(msg string, err error) *ClientError {
	return &ClientError{msg, err}
}

func (ce *ClientError) Error() string {
	return fmt.Sprintf("Audit client error with message [%s] with error [%v]", ce.message, ce.err)
}

func (ce *ClientError) Unwrap() error {
	return ce.err
}
