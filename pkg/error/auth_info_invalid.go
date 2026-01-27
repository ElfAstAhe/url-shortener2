package error

import "fmt"

type AuthInfoInvalidError struct {
	Msg string
	err error
}

var AuthInfoInvalidErr *AuthInfoInvalidError

func NewAuthInfoInvalidError(msg string, err error) *AuthInfoInvalidError {
	return &AuthInfoInvalidError{
		Msg: msg,
		err: err,
	}
}

func (e *AuthInfoInvalidError) Error() string {
	return fmt.Sprintf("auth info invalid with message [%s]", e.Msg)
}

func (e *AuthInfoInvalidError) Unwrap() error {
	return e.err
}
