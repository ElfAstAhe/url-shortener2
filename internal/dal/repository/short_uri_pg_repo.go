package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config/db"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/repository"
	apperrs "github.com/ElfAstAhe/url-shortener2/internal/error"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

const (
	getShortURISQL           string = "select id, original_url, key from short_uris where id = $1"
	getShortURIByKeySQL      string = "select id, original_url, key from short_uris where key = $1"
	getShortURIByKeyUserSQL  string = `select su.id, su.original_url, su.key, suu.deleted from short_uris su inner join short_uri_users suu on suu.short_uri_id = su.id and suu.user_id = $2 where su.key = $1`
	createShortURISQL        string = "insert into short_uris(id, original_url, key) values ($1, $2, $3)"
	listShortURIAllByUserSQL string = `select
    s.id,
    s.original_url,
    s.key
from
    short_uris s
    inner join short_uri_users su
        on
            su.user_id = $1
        and su.short_uri_id = s.id`
	listShortURIIdsByKeysSQL string = `select su.id from short_uris su where su.key = any($1)`
	getShortURICountSQL      string = `select count(1) as cnt from short_uris`
)

type ShortURIPgRepo struct {
	db       db.DB
	userRepo repository.ShortURIUserRepository
}

func NewShortURIPgRepo(appDB db.DB, userRepo repository.ShortURIUserRepository) (*ShortURIPgRepo, error) {
	if appDB == nil {
		return nil, errs.NewAppInvalidArgumentError("appDB", "nil")
	}

	return &ShortURIPgRepo{
		db:       appDB,
		userRepo: userRepo,
	}, nil
}

// ShortURIRepository

func (pgs *ShortURIPgRepo) Get(ctx context.Context, id string) (*model.ShortURI, error) {
	row := pgs.db.GetDB().QueryRowContext(ctx, getShortURISQL, id)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = model.ShortURI{
		OriginalURL: &model.CustomURL{},
	}

	// id, original_url, key
	err := row.Scan(&result.ID, result.OriginalURL, &result.Key)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pgs *ShortURIPgRepo) GetByKey(ctx context.Context, key string) (*model.ShortURI, error) {
	row := pgs.db.GetDB().QueryRowContext(ctx, getShortURIByKeySQL, key)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = model.ShortURI{
		OriginalURL: &model.CustomURL{},
	}

	// id, original_url, key, create_user, created, update_user, updated
	err := row.Scan(&result.ID, result.OriginalURL, &result.Key)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pgs *ShortURIPgRepo) GetByKeyUser(ctx context.Context, userID string, key string) (*model.ShortURI, error) {
	if userID == "" {
		return nil, nil
	}
	if key == "" {
		return nil, nil
	}

	row := pgs.db.GetDB().QueryRowContext(ctx, getShortURIByKeyUserSQL, key, userID)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = model.ShortURI{
		OriginalURL: &model.CustomURL{},
	}
	var deleted = false
	err := row.Scan(&result.ID, &result.OriginalURL, &result.Key, &deleted)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if deleted {
		return nil, apperrs.NewDalSoftRemovedError("short_uri", nil)
	}

	return &result, nil
}

func (pgs *ShortURIPgRepo) Create(ctx context.Context, userID string, entity *model.ShortURI) (*model.ShortURI, error) {
	if err := model.ValidateShortURI(entity); err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, errs.NewAuthInfoAbsentError("userID empty", nil)
	}

	find, err := pgs.GetByKey(ctx, entity.Key)
	if err != nil {
		return nil, err
	}
	if find != nil {
		err := pgs.addUser(ctx, nil, find.ID, userID)
		if err != nil && errors.As(err, &errs.ModelAlreadyExistsErr) {
			return find, err
		} else if err != nil {
			return nil, err
		}
		return find, errors.New("shortURI already exists")
	}

	tx, err := pgs.db.GetDB().Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()

			return
		}

		tx.Commit()
	}()

	stmt, err := tx.PrepareContext(ctx, createShortURISQL)
	if err != nil {
		return nil, err
	}
	defer utils.CloseOnly(stmt)
	stmtSU, err := tx.PrepareContext(ctx, createShortURIUserSQL)
	if err != nil {
		return nil, err
	}
	defer utils.CloseOnly(stmtSU)

	// id, original_url, key
	res, err := pgs.internalCreate(ctx, stmt, entity)
	if err != nil {
		return nil, err
	}
	if err := pgs.addUser(ctx, stmtSU, res.ID, userID); err != nil {
		return nil, err
	}

	return res, nil
}

// BatchCreate is creation a batch data in transaction
func (pgs *ShortURIPgRepo) BatchCreate(ctx context.Context, userID string, batch map[string]*model.ShortURI) (map[string]*model.ShortURI, error) {
	if len(batch) == 0 {
		return batch, nil
	}

	for _, entity := range batch {
		if err := model.ValidateShortURI(entity); err != nil {
			return nil, fmt.Errorf("batch validation, invalid entity: [%v] with error [%v]", entity, err)
		}
	}
	if userID == "" {
		return nil, errs.NewAuthInfoAbsentError("short uri batch create", nil)
	}

	tx, err := pgs.db.GetDB().Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()

			return
		}

		tx.Commit()
	}()

	stmt, err := tx.Prepare(createShortURISQL)
	if err != nil {
		return nil, err
	}
	defer utils.CloseOnly(stmt)
	stmtSU, err := tx.PrepareContext(ctx, createShortURIUserSQL)
	if err != nil {
		return nil, err
	}
	defer utils.CloseOnly(stmtSU)

	res := make(map[string]*model.ShortURI)
	for correlation, entity := range batch {
		find, err := pgs.GetByKey(ctx, entity.Key)
		if err != nil {
			return nil, err
		}
		if find != nil {
			if err := pgs.addUser(ctx, stmtSU, find.ID, userID); err != nil {
				return nil, err
			}
			res[correlation] = find

			continue
		}

		saved, err := pgs.internalCreate(ctx, stmt, entity)
		if err != nil {
			return nil, err
		}
		if err := pgs.addUser(ctx, stmtSU, saved.ID, userID); err != nil {
			return nil, err
		}

		res[correlation] = saved
	}

	return res, nil
}

func (pgs *ShortURIPgRepo) ListAllByKeys(ctx context.Context, keys []string) ([]*model.ShortURI, error) {
	//TODO implement me
	panic("implement me")
}

func (pgs *ShortURIPgRepo) ListAllByUser(ctx context.Context, userID string) ([]*model.ShortURI, error) {
	res := make([]*model.ShortURI, 0)
	if userID == "" {
		return res, nil
	}
	rows, err := pgs.db.GetDB().QueryContext(ctx, listShortURIAllByUserSQL, userID)
	if err != nil {
		return nil, err
	}
	defer utils.CloseOnly(rows)
	for rows.Next() {
		var result = model.ShortURI{
			OriginalURL: &model.CustomURL{},
		}

		err := rows.Scan(&result.ID, result.OriginalURL, &result.Key)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return res, nil
		} else if err != nil {
			return nil, err
		}

		res = append(res, &result)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}

func (pgs *ShortURIPgRepo) Delete(ctx context.Context, ID string, userID string) error {
	return pgs.userRepo.DeleteByUnique(ctx, userID, ID)
}

func (pgs *ShortURIPgRepo) BatchDeleteByKeys(ctx context.Context, userID string, keys []string) error {
	if userID == "" || len(keys) == 0 {
		return nil
	}

	ids, err := pgs.listIdsByKeys(ctx, keys)
	if err != nil {
		return err
	}

	return pgs.userRepo.DeleteAllByUnique(ctx, userID, ids)
}

func (pgs *ShortURIPgRepo) Count(ctx context.Context) (int, error) {
	row := pgs.db.GetDB().QueryRowContext(ctx, getShortURICountSQL)
	var res int
	err := row.Scan(&res)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return res, nil
}

func (pgs *ShortURIPgRepo) UniqueCount(ctx context.Context) (int, int, error) {
	eg, egCtx := errgroup.WithContext(ctx)
	var urlCount, userCount int
	eg.Go(func() error {
		var urlErr error
		urlCount, urlErr = pgs.Count(egCtx)

		return urlErr
	})
	eg.Go(func() error {
		var userErr error
		userCount, userErr = pgs.userRepo.UniqueCount(egCtx)

		return userErr
	})

	return urlCount, userCount, eg.Wait()
}

func (pgs *ShortURIPgRepo) internalCreate(ctx context.Context, preparedSQL *sql.Stmt, entity *model.ShortURI) (*model.ShortURI, error) {
	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	entity.ID = newID.String()

	_, err = preparedSQL.ExecContext(ctx, entity.ID, entity.OriginalURL, entity.Key)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (pgs *ShortURIPgRepo) addUser(ctx context.Context, stmt *sql.Stmt, id string, userID string) error {
	user, err := model.NewShortURIUser(id, userID)
	if err != nil {
		return err
	}
	if stmt != nil {
		if _, err := pgs.userRepo.CreateStmt(ctx, stmt, user); err != nil {
			return err
		}
	} else {
		if _, err := pgs.userRepo.Create(ctx, user); err != nil {
			return err
		}
	}

	return nil
}

func (pgs *ShortURIPgRepo) listIdsByKeys(ctx context.Context, keys []string) ([]string, error) {
	res := make([]string, 0)
	if len(keys) == 0 {
		return res, nil
	}

	rows, err := pgs.db.GetDB().QueryContext(ctx, listShortURIIdsByKeysSQL, keys)
	if err != nil {
		return nil, err
	}
	defer utils.CloseOnly(rows)

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return res, nil
		} else if err != nil {
			return nil, err
		}

		res = append(res, id)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}
