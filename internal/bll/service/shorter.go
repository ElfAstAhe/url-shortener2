package service

import (
	"context"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
)

// Shorter app service
type Shorter interface {
	// GetURL return full URL
	GetURL(ctx context.Context, key string) (string, error)
	GetURLUser(ctx context.Context, userID string, key string) (string, error)

	// Store URL and return short key
	Store(ctx context.Context, userID string, url string) (string, error)

	// BatchStore URLs and return correlation shorts
	BatchStore(ctx context.Context, userID string, source model.CorrelationUrls) (model.CorrelationShorts, error)

	// GetAllUserShorts return all user shorten urls
	GetAllUserShorts(ctx context.Context, userID string) (model.UserShorts, error)

	// BatchDelete remove short uris by user id
	BatchDelete(ctx context.Context, userID string, data model.UserBatchDeletes) error

	// InternalStats show internal data statistic info
	// returns
	//  total url count
	//  total unique user count
	//  error
	InternalStats(ctx context.Context) (int, int, error)
}
