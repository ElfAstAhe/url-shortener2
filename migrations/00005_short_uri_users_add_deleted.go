package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

const alterTableShortURIUsersAddDeleted = `alter table if exists short_uri_users add column if not exists deleted bool null default false;`

const alterTableShortURIUsersDropDeleted = `alter table if exists short_uri_users drop column if exists deleted;`

func init() {
	goose.AddMigrationNoTxContext(up00005, down00005)
}

func up00005(ctx context.Context, db *sql.DB) error {
	return alterTableShortURIUsersAddColumnDeleted(ctx, db)
}

func down00005(ctx context.Context, db *sql.DB) error {
	return alterTableShortURIUsersDropColumnDeleted(ctx, db)
}

func alterTableShortURIUsersAddColumnDeleted(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, alterTableShortURIUsersAddDeleted)

	return err
}

func alterTableShortURIUsersDropColumnDeleted(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, alterTableShortURIUsersDropDeleted)

	return err
}
