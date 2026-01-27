package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

const createTableShortURIUsers = `CREATE TABLE IF NOT EXISTS short_uri_users (
    id varchar(50) not null,
    short_uri_id varchar(50) not null,
    user_id varchar(50) not null,
    constraint short_uri_users_pk PRIMARY KEY (id)
);`

const dropTableShortURIUsers = `drop table if exists short_uri_users cascade;`

const createIndexShortURIUsersUser = `create index if not exists short_uri_users_user_su_idx on short_uri_users(user_id asc, short_uri_id asc);`

const dropIndexShortURIUsersUser = `drop index if exists short_uri_users_user_su_idx cascade;`

func init() {
	goose.AddMigrationNoTxContext(up00004, down00004)
}

func up00004(ctx context.Context, db *sql.DB) error {
	if err := upCreateTableShortURIUsers(ctx, db); err != nil {
		return err
	}

	// create index
	return upCreateIndexShortURIUsers(ctx, db)
}

func down00004(ctx context.Context, db *sql.DB) error {
	// drop index
	if err := downDropIndexShortURIUsers(ctx, db); err != nil {
		return err
	}

	// drop table
	return downDropTableShortURIUsers(ctx, db)
}

func upCreateTableShortURIUsers(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createTableShortURIUsers)
	if err != nil {
		return err
	}

	return nil
}

func downDropTableShortURIUsers(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, dropTableShortURIUsers)
	if err != nil {
		return err
	}

	return nil
}

func upCreateIndexShortURIUsers(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createIndexShortURIUsersUser)
	if err != nil {
		return err
	}

	return nil
}

func downDropIndexShortURIUsers(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, dropIndexShortURIUsersUser)
	if err != nil {
		return err
	}

	return nil
}
