package facade

import (
	"context"
)

type ToolFacade interface {
	Ping(ctx context.Context) error
}
