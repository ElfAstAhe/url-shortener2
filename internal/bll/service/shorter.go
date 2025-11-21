package service

import "context"

type CorrelationUrls map[string]string

type CorrelationShorts map[string]string

type UserShorts map[string]string

// UserBatchDeletes key is user id values is short uri keys
type UserBatchDeletes []string

// Shorter app service
type Shorter interface {
	Close() error

	// GetURL return full URL
	GetURL(ctx context.Context, key string) (string, error)

	// Store URL and return short key
	Store(ctx context.Context, url string) (string, error)

	// BatchStore URLs and return correlation shorts
	BatchStore(ctx context.Context, source CorrelationUrls) (CorrelationShorts, error)

	// GetAllUserShorts return all user shorten urls
	GetAllUserShorts(ctx context.Context, userID string) (UserShorts, error)

	// BatchDelete remove short uris by user id
	BatchDelete(ctx context.Context, data UserBatchDeletes) error
}
