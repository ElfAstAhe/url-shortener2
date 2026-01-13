package facade

import (
	"context"
	"io"
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/internal/ep/dto"
)

type ShortenFacade interface {
	GetURL(ctx context.Context, req *http.Request) (string, error)
	CreateURL(ctx context.Context, income io.Reader) (*dto.ShortenCreateResponse, error)
	BatchCreateURL(ctx context.Context, income io.Reader) ([]*dto.ShortenBatchResponseItem, error)
	Store(ctx context.Context, income io.Reader) (string, error)
}
