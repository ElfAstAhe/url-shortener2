package facade

import (
	"context"
	"net/http"
)

type ShortenFacade interface {
	GetURL(ctx context.Context, req *http.Request) (string, error)
}
