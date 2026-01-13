package error

import "fmt"

type BllConflictError struct {
	Key string
	err error
}

var BllConflictErr *BllInvalidOriginalURLError

func NewBllConflictError(key string) *BllConflictError {
	return NewBllConflictErrorEx(key, nil)
}

func NewBllConflictErrorEx(key string, err error) *BllConflictError {
	return &BllConflictError{
		Key: key,
		err: err,
	}
}

func (e *BllConflictError) Error() string {
	return fmt.Sprintf("Bll conflict error with key [%s]", e.Key)
}

func (e *BllConflictError) Unwrap() error {
	return e.err
}
