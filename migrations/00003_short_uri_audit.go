package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

const createTableShortURIAudit = `CREATE TABLE IF NOT EXISTS short_uri_audit (
    id varchar(50) not null,
    short_uri_id varchar(50) not null,
    user_id varchar(50) not null,
    "user" varchar(100) null,
    date timestamptz not null default now(),
    operation varchar(50) not null,
    constraint short_uri_audit_pk PRIMARY KEY (id)
);`

const dropTableShortURIAudit = `drop table if exists short_uri_audit cascade;`

const createIndexShortURIAuditUser = `create index if not exists short_uri_audit_user_su_idx on short_uri_audit(user_id asc, short_uri_id asc);`

const dropIndexShortURIAuditUser = `drop index if exists short_uri_audit_user_su cascade;`

func init() {
	goose.AddMigrationNoTxContext(up00003, down00003)
}

func up00003(ctx context.Context, db *sql.DB) error {
	if err := upCreateTableShortURIAudit(ctx, db); err != nil {
		return err
	}

	// create index
	return upCreateIndexShortURIAudit(ctx, db)
}

func down00003(ctx context.Context, db *sql.DB) error {
	// drop index
	if err := downDropIndexShortURIAudit(ctx, db); err != nil {
		return err
	}

	// drop table
	return downDropTableShortURIAudit(ctx, db)
}

func upCreateTableShortURIAudit(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createTableShortURIAudit)
	if err != nil {
		return err
	}

	return nil
}

func downDropTableShortURIAudit(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, dropTableShortURIAudit)
	if err != nil {
		return err
	}

	return nil
}

func upCreateIndexShortURIAudit(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createIndexShortURIAuditUser)
	if err != nil {
		return err
	}

	return nil
}

func downDropIndexShortURIAudit(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, dropIndexShortURIAuditUser)
	if err != nil {
		return err
	}

	return nil
}
