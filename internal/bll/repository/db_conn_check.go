package repository

import (
	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	"github.com/ElfAstAhe/url-shortener2/internal/app/config/db"
	"github.com/ElfAstAhe/url-shortener2/internal/dal/reppository"
)

type DBConnCheckRepository interface {
	CheckDBConn() error
	Close() error
}

func NewDBConnCheckRepository(appDb db.DB) (DBConnCheckRepository, error) {
	if appDb.GetDBKind() == config.DBKindPostgres {
		return repository.NewDBConnCheckPgRepo(appDb)
	}

	return repository.NewDBConnCheckImMemRepo()
}
