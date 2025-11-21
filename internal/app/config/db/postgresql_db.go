package db

import (
	"database/sql"
	"time"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgresqlDB struct {
	DB     *sql.DB
	DBKind string
	Dsn    string
}

func newPostgresqlDB(dsn string) (*postgresqlDB, error) {
	pg, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	pg.SetMaxOpenConns(20)
	pg.SetMaxIdleConns(5)
	pg.SetConnMaxIdleTime(60 * time.Second)

	err = pg.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresqlDB{
		DB:     pg,
		DBKind: config.DBKindPostgres,
		Dsn:    dsn,
	}, nil
}

// Closer

func (pdb *postgresqlDB) Close() error {
	return pdb.DB.Close()
}

// =============

// DB

func (pdb *postgresqlDB) GetDB() *sql.DB {
	return pdb.DB
}

func (pdb *postgresqlDB) GetDBKind() string {
	return pdb.DBKind
}

func (pdb *postgresqlDB) GetDsn() string {
	return pdb.Dsn
}

// =============
