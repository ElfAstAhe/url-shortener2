package facade

import (
	"context"
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/internal/ep/dto"
)

type ShortenFacade interface {
	GetURL(ctx context.Context, req *http.Request) (string, error)
	CreateURL(ctx context.Context, req *http.Request) (*dto.ShortenCreateResponse, error)
	BatchCreateURL(ctx context.Context, r *http.Request) ([]*dto.ShortenBatchResponseItem, error)
	Store(ctx context.Context, r *http.Request) (string, error)
}
