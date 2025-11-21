package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

const createIndexShortURIs = `create index if not exists short_uris_key_idx on short_uris(key asc)`

const dropIndexShortURIs = `drop index if exists short_uris_key_idx`

func init() {
	goose.AddMigrationNoTxContext(Up00002, Down00002)
}

func Up00002(ctx context.Context, db *sql.DB) error {
	return upCreateIndexShortURIsKeyIdx(ctx, db)
}

func Down00002(ctx context.Context, db *sql.DB) error {
	return downDropIndexShortURIsKeyIdx(ctx, db)
}

func upCreateIndexShortURIsKeyIdx(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createIndexShortURIs)
	if err != nil {
		return err
	}

	return nil
}

func downDropIndexShortURIsKeyIdx(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, dropIndexShortURIs)
	if err != nil {
		return err
	}

	return nil
}
