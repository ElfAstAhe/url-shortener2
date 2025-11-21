package repository

import (
	"context"
	"time"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config/db"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

type DBConnCheckPgRepo struct {
	DB db.DB
}

func NewDBConnCheckPgRepo(appDB db.DB) (*DBConnCheckPgRepo, error) {
	if appDB == nil {
		return nil, errs.NewAppInvalidArgumentError("appDB", "nil")
	}

	return &DBConnCheckPgRepo{
		DB: appDB,
	}, nil
}

// Closer

func (pgr *DBConnCheckPgRepo) Close() error {
	return db.CloseDB(pgr.DB)
}

// ========

// DBConn

func (pgr *DBConnCheckPgRepo) CheckDBConn() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return pgr.DB.GetDB().PingContext(ctx)
}

// ========
