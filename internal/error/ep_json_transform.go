package error

import "fmt"

type EpJSONTransformError struct {
	targetType string
	err        error
}

var EpJSONTransformErr *EpJSONTransformError

func NewEpJSONTransformError(targetType string, err error) *EpJSONTransformError {
	return &EpJSONTransformError{
		targetType: targetType,
		err:        err,
	}
}

func (jte *EpJSONTransformError) Error() string {
	return fmt.Sprintf("error transform json data into target type [%s] with error [%v]", jte.targetType, jte.err)
}

func (jte *EpJSONTransformError) Unwrap() error {
	return jte.err
}
