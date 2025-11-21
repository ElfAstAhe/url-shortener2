package repository

import (
	errs "github.com/ElfAstAhe/url-shortener2/internal/error"
)

type DBConnCheckImMemRepo struct {
}

func NewDBConnCheckImMemRepo() (*DBConnCheckImMemRepo, error) {
	return &DBConnCheckImMemRepo{}, nil
}

func (D *DBConnCheckImMemRepo) CheckDBConn() error {
	return errs.NewDalDBConnCheckError("IN_MEMORY DB, where is no connection")
}

func (D *DBConnCheckImMemRepo) Close() error {
	return nil
}
