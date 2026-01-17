package error

import "fmt"

type AuthCookieInvalidError struct {
	Msg string
	err error
}

var AuthCookieInvalidErr *AuthCookieInvalidError

func NewAuthCookieInvalidError(msg string, err error) *AuthCookieInvalidError {
	return &AuthCookieInvalidError{
		Msg: msg,
		err: err,
	}
}

func (e *AuthCookieInvalidError) Error() string {
	return fmt.Sprintf("auth cookie invalid with message [%s]", e.Msg)
}

func (e *AuthCookieInvalidError) Unwrap() error {
	return e.err
}
