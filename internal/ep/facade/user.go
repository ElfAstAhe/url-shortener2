package facade

import (
	"context"
	"io"

	"github.com/ElfAstAhe/url-shortener2/internal/ep/dto"
)

type UserFacade interface {
	BatchDelete(ctx context.Context, rawData io.Reader) error
	ListAllShortens(ctx context.Context) ([]*dto.UserShorten, error)
}
