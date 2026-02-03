package facade

import (
	"context"
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/internal/ep/dto"
)

type ToolFacade interface {
	Ping(ctx context.Context) error
	InternalStats(ctx context.Context, header *http.Header) (*dto.InternalStatsDto, error)
}
