package db

import (
	"database/sql"
	"io"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
)

type DB interface {
	GetDB() *sql.DB
	GetDBKind() string
	GetDsn() string
}

func NewDB(kind string, dsn string) (DB, error) {
	if kind == config.DBKindPostgres {
		return newPostgresqlDB(dsn)
	}

	return newInMemoryDB()
}

func CloseDB(db DB) error {
	if closer, ok := db.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
