package error

import "fmt"

type BllInvalidOriginalURLError struct {
	OriginalURL string
}

var BllInvalidOriginalURLErr *BllInvalidOriginalURLError

func NewBllInvalidOriginalURLError(originalURL string) *BllInvalidOriginalURLError {
	return &BllInvalidOriginalURLError{OriginalURL: originalURL}
}

func (i *BllInvalidOriginalURLError) Error() string {
	return fmt.Sprintf("originalURL ['%s'] is invalid", i.OriginalURL)
}
