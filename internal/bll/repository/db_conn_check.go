package repository

import "context"

type DBConnCheckRepository interface {
	CheckDBConn(ctx context.Context) error
	Close() error
}
