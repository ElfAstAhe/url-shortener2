package error

import (
	"fmt"
)

type DalSoftRemovedError struct {
	Entity string
	Err    error
}

var DalSoftRemovedErr *DalSoftRemovedError

func NewDalSoftRemovedError(entity string, err error) *DalSoftRemovedError {
	return &DalSoftRemovedError{
		Entity: entity,
		Err:    err,
	}
}

func (e *DalSoftRemovedError) Error() string {
	return fmt.Sprintf("entity [%s] soft removed", e.Entity)
}

func (e *DalSoftRemovedError) Unwrap() error { return e.Err }
