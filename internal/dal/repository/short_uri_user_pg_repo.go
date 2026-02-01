package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config/db"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

type ShortURIUserPgRepo struct {
	db db.DB
}

const (
	getShortURIUserSQL                 string = `select id, short_uri_id, user_id, deleted from short_uri_users where id = $1`
	getShortURIUserByUniqueSQL         string = `select id, short_uri_id, user_id, deleted from short_uri_users where user_id = $1 and short_uri_id = $2`
	listShortURIUserAllByUserSQL       string = `select id, short_uri_id, user_id, deleted from short_uri_users where user_id = $1`
	listShortURIUserAllByShortURISQL   string = `select id, short_uri_id, user_id, deleted from short_uri_users where short_uri_id = $1`
	createShortURIUserSQL              string = `insert into short_uri_users(id, short_uri_id, user_id, deleted) values($1, $2, $3, $4)`
	changeShortURIUserSQL              string = `update short_uri_users set short_uri_id=$2, user_id=$3, deleted = $4 where id = $1`
	deleteShortURIUserSQL              string = `update short_uri_users set deleted = true where id = $1`
	deleteShortURIUserByUniqueSQL      string = `update short_uri_users set deleted = true where user_id = $1 and short_uri_id = $2`
	deleteAllShortURIUserByUniqueSQL   string = `update short_uri_users set deleted = true where user_id = $1 and short_uri_id = any($2)`
	deleteAllShortURIUserByUserSQL     string = `update short_uri_users set deleted = true where user_id = $1`
	deleteAllShortURIUserByShortURISQL string = `update short_uri_users set deleted = true where short_uri_id = $1`
	removeShortURIUserSQL              string = `delete from short_uri_users where id = $1`
	removeShortURIUserByUniqueSQL      string = `delete from short_uri_users where user_id = $1 and short_uri_id = $2`
	removeAllShortURIUserByUniqueSQL   string = `delete from short_uri_users where user_id = $1 and short_uri_id = any($2)`
	removeAllShortURIUserByUserSQL     string = `delete from short_uri_users where user_id = $1`
	removeAllShortURIUserByShortURISQL string = `delete from short_uri_users where short_uri_id = $1`
	getShortURIUserCountSQL            string = `select count(id) from short_uri_users`
	getShortURIUserUniqueCountSQL      string = `select count(1) as ucnt from (select distinct user_id from short_uri_users)`
)

func NewShortURIUserPgRepo(appDB db.DB) (*ShortURIUserPgRepo, error) {
	if appDB == nil {
		return nil, errs.NewAppInvalidArgumentError("appDB", "nil")
	}

	return &ShortURIUserPgRepo{
		db: appDB,
	}, nil
}

func (pgsu *ShortURIUserPgRepo) Get(ctx context.Context, ID string) (*model.ShortURIUser, error) {
	row := pgsu.db.GetDB().QueryRow(getShortURIUserSQL, ID)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = model.ShortURIUser{}
	// id, short_uri_id, user_id, deleted
	err := row.Scan(&result.ID, &result.ShortURIID, &result.UserID, &result.Deleted)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pgsu *ShortURIUserPgRepo) GetByUnique(ctx context.Context, userID string, shortURIID string) (*model.ShortURIUser, error) {
	row := pgsu.db.GetDB().QueryRow(getShortURIUserByUniqueSQL, userID, shortURIID)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = model.ShortURIUser{}
	// id, short_uri_id, user_id, deleted
	err := row.Scan(&result.ID, &result.ShortURIID, &result.UserID, &result.Deleted)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pgsu *ShortURIUserPgRepo) ListAllByUser(ctx context.Context, userID string) ([]*model.ShortURIUser, error) {
	res := make([]*model.ShortURIUser, 0)
	rows, err := pgsu.db.GetDB().Query(listShortURIUserAllByUserSQL, userID)
	if err != nil {
		return res, err
	}
	defer utils.CloseOnly(rows)

	for rows.Next() {
		var result = model.ShortURIUser{}
		// id, short_uri_id, user_id, deleted
		err := rows.Scan(&result.ID, &result.ShortURIID, &result.UserID, &result.Deleted)
		if err != nil {
			return nil, err
		}

		res = append(res, &result)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}

func (pgsu *ShortURIUserPgRepo) ListAllByShortURI(ctx context.Context, shortURIID string) ([]*model.ShortURIUser, error) {
	res := make([]*model.ShortURIUser, 0)
	rows, err := pgsu.db.GetDB().Query(listShortURIUserAllByShortURISQL, shortURIID)
	if err != nil {
		return res, err
	}
	defer utils.CloseOnly(rows)

	for rows.Next() {
		var result = model.ShortURIUser{}
		// id, short_uri_id, user_id, deleted
		err := rows.Scan(&result.ID, &result.ShortURIID, &result.UserID, &result.Deleted)
		if err != nil {
			return nil, err
		}

		res = append(res, &result)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}

func (pgsu *ShortURIUserPgRepo) Create(ctx context.Context, entity *model.ShortURIUser) (*model.ShortURIUser, error) {
	if err := model.ValidateShortURIUser(entity); err != nil {
		return nil, err
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, createShortURIUserSQL)
	if err != nil {
		return nil, err
	}
	defer utils.CloseOnly(stmt)

	res, err := pgsu.CreateStmt(ctx, stmt, entity)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (pgsu *ShortURIUserPgRepo) CreateStmt(ctx context.Context, stmt *sql.Stmt, entity *model.ShortURIUser) (*model.ShortURIUser, error) {
	if err := model.ValidateShortURIUser(entity); err != nil {
		return nil, err
	}

	find, err := pgsu.GetByUnique(ctx, entity.UserID, entity.ShortURIID)
	if err != nil {
		return nil, err
	}
	if find != nil {
		return find, errs.NewModelAlreadyExistsError("short_uri_user", entity.ID)
	}

	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	entity.ID = newID.String()

	_, err = stmt.ExecContext(ctx, entity.ID, entity.ShortURIID, entity.UserID, entity.Deleted)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (pgsu *ShortURIUserPgRepo) Change(ctx context.Context, entity *model.ShortURIUser) (*model.ShortURIUser, error) {
	if err := model.ValidateShortURIUser(entity); err != nil {
		return nil, errs.NewModelValidationError("short_uri_user", "", err)
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, changeShortURIUserSQL)
	if err != nil {
		return nil, err
	}
	defer utils.CloseOnly(stmt)

	return pgsu.ChangeStmt(ctx, stmt, entity)
}

func (pgsu *ShortURIUserPgRepo) Count(ctx context.Context) (int, error) {
	row := pgsu.db.GetDB().QueryRowContext(ctx, getShortURIUserCountSQL)
	var count int
	err := row.Scan(&count)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return count, nil
}

func (pgsu *ShortURIUserPgRepo) UniqueCount(ctx context.Context) (int, error) {
	row := pgsu.db.GetDB().QueryRowContext(ctx, getShortURIUserUniqueCountSQL)
	var count int
	err := row.Scan(&count)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return count, nil
}

func (pgsu *ShortURIUserPgRepo) ChangeStmt(ctx context.Context, stmt *sql.Stmt, entity *model.ShortURIUser) (*model.ShortURIUser, error) {
	if err := model.ValidateShortURIUser(entity); err != nil {
		return nil, errs.NewModelValidationError("short_uri_user", "", err)
	}

	_, err := stmt.ExecContext(ctx, entity.ID, entity.ShortURIID, entity.UserID, entity.Deleted)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (pgsu *ShortURIUserPgRepo) Delete(ctx context.Context, ID string) error {
	if ID == "" {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, deleteShortURIUserSQL)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(stmt)

	return pgsu.DeleteStmt(ctx, stmt, ID)
}

func (pgsu *ShortURIUserPgRepo) DeleteStmt(ctx context.Context, stmt *sql.Stmt, ID string) error {
	if ID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *ShortURIUserPgRepo) DeleteByUnique(ctx context.Context, userID string, shortURIID string) error {
	if userID == "" || shortURIID == "" {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, deleteShortURIUserByUniqueSQL)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(stmt)

	return pgsu.DeleteByUniqueStmt(ctx, stmt, userID, shortURIID)
}

func (pgsu *ShortURIUserPgRepo) DeleteByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIID string) error {
	if userID == "" || shortURIID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, userID, shortURIID)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *ShortURIUserPgRepo) DeleteAllByUnique(ctx context.Context, userID string, shortURIIDs []string) error {
	if userID == "" || len(shortURIIDs) == 0 {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, deleteAllShortURIUserByUniqueSQL)
	if err != nil {
		return err
	}

	return pgsu.DeleteAllByUniqueStmt(ctx, stmt, userID, shortURIIDs)
}

func (pgsu *ShortURIUserPgRepo) DeleteAllByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIIDs []string) error {
	if userID == "" || len(shortURIIDs) == 0 {
		return nil
	}

	_, err := stmt.ExecContext(ctx, userID, shortURIIDs)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *ShortURIUserPgRepo) DeleteAllByUser(ctx context.Context, userID string) error {
	if userID == "" {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, deleteAllShortURIUserByUserSQL)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(stmt)

	return pgsu.DeleteAllByUserStmt(ctx, stmt, userID)
}

func (pgsu *ShortURIUserPgRepo) DeleteAllByUserStmt(ctx context.Context, stmt *sql.Stmt, userID string) error {
	if userID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *ShortURIUserPgRepo) DeleteAllByShortURI(ctx context.Context, shortURIID string) error {
	if shortURIID == "" {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, deleteAllShortURIUserByShortURISQL)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(stmt)

	return pgsu.DeleteAllByShortURIStmt(ctx, stmt, shortURIID)
}

func (pgsu *ShortURIUserPgRepo) DeleteAllByShortURIStmt(ctx context.Context, stmt *sql.Stmt, shortURIID string) error {
	if shortURIID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, shortURIID)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *ShortURIUserPgRepo) Remove(ctx context.Context, ID string) error {
	if ID == "" {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, removeShortURIUserSQL)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(stmt)

	return pgsu.RemoveStmt(ctx, stmt, ID)
}

func (pgsu *ShortURIUserPgRepo) RemoveStmt(ctx context.Context, stmt *sql.Stmt, ID string) error {
	if ID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *ShortURIUserPgRepo) RemoveByUnique(ctx context.Context, userID string, shortURIID string) error {
	if userID == "" || shortURIID == "" {
		return nil
	}

	stmp, err := pgsu.db.GetDB().PrepareContext(ctx, removeShortURIUserByUniqueSQL)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(stmp)

	return pgsu.RemoveByUniqueStmt(ctx, stmp, userID, shortURIID)
}

func (pgsu *ShortURIUserPgRepo) RemoveByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIID string) error {
	if userID == "" || shortURIID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, userID, shortURIID)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *ShortURIUserPgRepo) RemoveAllByUnique(ctx context.Context, userID string, shortURIIDs []string) error {
	if userID == "" || len(shortURIIDs) == 0 {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, removeAllShortURIUserByUniqueSQL)
	if err != nil {
		return err
	}

	return pgsu.RemoveAllByUniqueStmt(ctx, stmt, userID, shortURIIDs)
}

func (pgsu *ShortURIUserPgRepo) RemoveAllByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIIds []string) error {
	if userID == "" || len(shortURIIds) == 0 {
		return nil
	}

	_, err := stmt.ExecContext(ctx, userID, shortURIIds)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *ShortURIUserPgRepo) RemoveAllByUser(ctx context.Context, userID string) error {
	if userID == "" {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, removeAllShortURIUserByUserSQL)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(stmt)

	return pgsu.RemoveAllByUserStmt(ctx, stmt, userID)
}

func (pgsu *ShortURIUserPgRepo) RemoveAllByUserStmt(ctx context.Context, stmt *sql.Stmt, userID string) error {
	if userID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (pgsu *ShortURIUserPgRepo) RemoveAllByShortURI(ctx context.Context, shortURIID string) error {
	if shortURIID == "" {
		return nil
	}

	stmt, err := pgsu.db.GetDB().PrepareContext(ctx, removeAllShortURIUserByShortURISQL)
	if err != nil {
		return err
	}
	defer utils.CloseOnly(stmt)

	return pgsu.RemoveAllByShortURIStmt(ctx, stmt, shortURIID)
}

func (pgsu *ShortURIUserPgRepo) RemoveAllByShortURIStmt(ctx context.Context, stmt *sql.Stmt, shortURIID string) error {
	if shortURIID == "" {
		return nil
	}

	_, err := stmt.ExecContext(ctx, shortURIID)
	if err != nil {
		return err
	}

	return nil
}
