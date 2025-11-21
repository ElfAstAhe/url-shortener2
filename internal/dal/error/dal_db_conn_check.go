package error

import "fmt"

type DalDBConnCheckError struct {
	message string
}

var DalDBConnCheckErr *DalDBConnCheckError

func NewDalDBConnCheckError(message string) *DalDBConnCheckError {
	return &DalDBConnCheckError{
		message: message,
	}
}

func (e *DalDBConnCheckError) Error() string {
	return fmt.Sprintf("error check db connection with message [%s]", e.message)
}
