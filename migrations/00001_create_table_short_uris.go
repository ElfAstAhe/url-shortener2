package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

const createTableShortURIs = `create table if not exists short_uris (
    id varchar(50) not null,
    original_url varchar(4096) not null,
    key varchar(50) not null,
    create_user varchar(50) null,
    created timestamptz null,
    update_user varchar(50) null,
    updated timestamptz null,
    constraint short_uris_pk primary key(id),
    constraint short_uris_uk unique(original_url)
);`

const dropTableShortURIs = `drop table if exists short_uris cascade;`

func init() {
	goose.AddMigrationNoTxContext(up00001, down00001)
}

func up00001(ctx context.Context, db *sql.DB) error {
	// create table short_uris
	return upCreateTableShortUris(ctx, db)
}

func down00001(ctx context.Context, db *sql.DB) error {
	return downDropTableShortUris(ctx, db)
}

func upCreateTableShortUris(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createTableShortURIs)
	if err != nil {
		return err
	}

	return nil
}

func downDropTableShortUris(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, dropTableShortURIs)
	if err != nil {
		return err
	}

	return nil
}
