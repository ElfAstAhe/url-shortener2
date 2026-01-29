package error

import "fmt"

type AppInvalidConfigError struct {
	Name    string
	Value   any
	Message string
	Err     error
	general bool
}

var AppInvalidConfigErr *AppInvalidConfigError

func NewAppInvalidConfigError(name string, value any, message string, err error) *AppInvalidConfigError {
	return &AppInvalidConfigError{
		Name:    name,
		Value:   value,
		Message: message,
		Err:     err,
		general: false,
	}
}

func NewAppInvalidConfigItemMsgError(name string, value any, message string) *AppInvalidConfigError {
	return NewAppInvalidConfigError(name, value, message, nil)
}

func NewAppInvalidConfigItemErr(name string, value any) *AppInvalidConfigError {
	return NewAppInvalidConfigItemMsgError(name, value, "")
}

func NewAppGeneralInvalidConfigError(message string, err error) *AppInvalidConfigError {
	return &AppInvalidConfigError{
		Name:    "General",
		Value:   nil,
		Message: message,
		Err:     err,
		general: true,
	}
}

func (apc *AppInvalidConfigError) Error() string {
	if apc.general {
		if apc.Err != nil {
			return fmt.Sprintf("general config error with message [%s] and source error [%v]", apc.Message, apc.Err)
		}

		return fmt.Sprintf("general config error with message [%s]", apc.Message)
	}

	return fmt.Sprintf("invalid config item with name [%s] with value [%v] with message [%s]", apc.Name, apc.Value, apc.Message)
}

func (apc *AppInvalidConfigError) Unwrap() error {
	return apc.Err
}
