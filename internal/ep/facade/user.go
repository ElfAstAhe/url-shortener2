package facade

import (
	"context"
	"io"
)

type UserFacade interface {
	BatchDelete(ctx context.Context, rawData io.Reader) error
}
