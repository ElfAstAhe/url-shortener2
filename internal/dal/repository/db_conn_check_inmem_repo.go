package repository

import (
	"context"
)

type DBConnCheckImMemRepo struct {
}

func NewDBConnCheckImMemRepo() (*DBConnCheckImMemRepo, error) {
	return &DBConnCheckImMemRepo{}, nil
}

func (D *DBConnCheckImMemRepo) CheckDBConn(ctx context.Context) error {
	return nil
}

func (D *DBConnCheckImMemRepo) Close() error {
	return nil
}
