package repository

import (
	"context"
	"database/sql"
	"maps"
	"slices"

	"github.com/google/uuid"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config/db"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

type ShortURIUserInMemRepo struct {
	cache db.InMemoryCache
}

func NewShortURIUserInMemRepo(appDB db.DB) (*ShortURIUserInMemRepo, error) {
	if cache, ok := appDB.(db.InMemoryCache); ok {
		return &ShortURIUserInMemRepo{
			cache: cache,
		}, nil
	}

	return nil, errs.NewAppInvalidConfigError("appDB", appDB, "db param does not implement InMemoryCache", nil)
}

func (imsu *ShortURIUserInMemRepo) Get(ctx context.Context, ID string) (*model.ShortURIUser, error) {
	if ID == "" {
		return nil, nil
	}

	imsu.cache.GetShortURIUserRWMutex().RLock()
	defer imsu.cache.GetShortURIUserRWMutex().RUnlock()
	return imsu.cache.GetShortURIUserCache()[ID], nil
}

func (imsu *ShortURIUserInMemRepo) GetByUnique(ctx context.Context, userID string, shortURIID string) (*model.ShortURIUser, error) {
	if userID == "" || shortURIID == "" {
		return nil, nil
	}

	imsu.cache.GetShortURIUserRWMutex().RLock()
	defer imsu.cache.GetShortURIUserRWMutex().RUnlock()
	for _, entity := range imsu.cache.GetShortURIUserCache() {
		if entity.UserID == userID && entity.ShortURIID == shortURIID {
			return entity, nil
		}
	}

	return nil, nil
}

func (imsu *ShortURIUserInMemRepo) ListAllByUser(ctx context.Context, userID string) ([]*model.ShortURIUser, error) {
	res := make([]*model.ShortURIUser, 0)
	if userID == "" || len(imsu.cache.GetShortURIUserCache()) == 0 {
		return res, nil
	}

	imsu.cache.GetShortURIUserRWMutex().RLock()
	defer imsu.cache.GetShortURIUserRWMutex().RUnlock()
	for _, entity := range imsu.cache.GetShortURIUserCache() {
		if entity.UserID == userID && !entity.Deleted {
			res = append(res, entity)
		}
	}

	return res, nil
}

func (imsu *ShortURIUserInMemRepo) ListAllByShortURI(ctx context.Context, shortURIID string) ([]*model.ShortURIUser, error) {
	res := make([]*model.ShortURIUser, 0)
	if shortURIID == "" || len(imsu.cache.GetShortURIUserCache()) == 0 {
		return res, nil
	}

	imsu.cache.GetShortURIUserRWMutex().RLock()
	defer imsu.cache.GetShortURIUserRWMutex().RUnlock()
	for _, entity := range imsu.cache.GetShortURIUserCache() {
		if entity.ShortURIID == shortURIID {
			res = append(res, entity)
		}
	}

	return res, nil
}

func (imsu *ShortURIUserInMemRepo) Create(ctx context.Context, entity *model.ShortURIUser) (*model.ShortURIUser, error) {
	if err := model.ValidateShortURIUser(entity); err != nil {
		return nil, errs.NewModelValidationError("short_uri_user", "", err)
	}

	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	entity.ID = newID.String()
	imsu.cache.GetShortURIUserRWMutex().Lock()
	defer imsu.cache.GetShortURIUserRWMutex().Unlock()
	imsu.cache.GetShortURIUserCache()[entity.ID] = entity

	return entity, nil
}

func (imsu *ShortURIUserInMemRepo) Change(ctx context.Context, entity *model.ShortURIUser) (*model.ShortURIUser, error) {
	if err := model.ValidateShortURIUser(entity); err != nil {
		return nil, errs.NewModelValidationError("short_uri_user", "", err)
	}
	find, err := imsu.Get(ctx, entity.ID)
	if err != nil {
		return nil, err
	}
	if find == nil {
		return nil, errs.NewModelNotExistsError("short_uri_user", entity.ID)
	}

	imsu.cache.GetShortURIUserRWMutex().Lock()
	defer imsu.cache.GetShortURIUserRWMutex().Unlock()
	imsu.cache.GetShortURIUserCache()[entity.ID] = entity

	return entity, nil
}

func (imsu *ShortURIUserInMemRepo) Delete(ctx context.Context, ID string) error {
	if ID == "" {
		return nil
	}

	find, err := imsu.Get(ctx, ID)
	if err != nil {
		return err
	}
	if find == nil {
		return nil
	}

	return imsu.delete(ctx, find)
}

func (imsu *ShortURIUserInMemRepo) DeleteByUnique(ctx context.Context, userID string, ID string) error {
	if userID == "" || ID == "" {
		return nil
	}

	find, err := imsu.GetByUnique(ctx, userID, ID)
	if err != nil {
		return err
	}
	if find == nil {
		return nil
	}

	return imsu.delete(ctx, find)
}

func (imsu *ShortURIUserInMemRepo) DeleteAllByUnique(ctx context.Context, userID string, shortURIIds []string) error {
	// ToDo: implement goroutine
	// ..

	if userID == "" || len(shortURIIds) == 0 {
		return nil
	}

	for _, shortURIID := range shortURIIds {
		err := imsu.DeleteByUnique(ctx, userID, shortURIID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (imsu *ShortURIUserInMemRepo) DeleteAllByUser(ctx context.Context, userID string) error {
	if userID == "" {
		return nil
	}

	toDelete, err := imsu.ListAllByUser(ctx, userID)
	if err != nil {
		return err
	}
	for _, entity := range toDelete {
		err := imsu.delete(ctx, entity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (imsu *ShortURIUserInMemRepo) DeleteAllByShortURI(ctx context.Context, shortURIID string) error {
	if shortURIID == "" {
		return nil
	}

	toDelete, err := imsu.ListAllByShortURI(ctx, shortURIID)
	if err != nil {
		return err
	}

	for _, entity := range toDelete {
		err := imsu.Delete(ctx, entity.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (imsu *ShortURIUserInMemRepo) Remove(ctx context.Context, ID string) error {
	if ID == "" {
		return nil
	}

	imsu.cache.GetShortURIUserRWMutex().Lock()
	defer imsu.cache.GetShortURIUserRWMutex().Unlock()
	delete(imsu.cache.GetShortURIUserCache(), ID)

	return nil
}

func (imsu *ShortURIUserInMemRepo) RemoveByUnique(ctx context.Context, userID string, shortURIID string) error {
	if userID == "" || shortURIID == "" {
		return nil
	}

	find, err := imsu.GetByUnique(ctx, userID, shortURIID)
	if err != nil {
		return err
	}
	if find == nil {
		return nil
	}

	return imsu.Remove(ctx, find.ID)
}

func (imsu *ShortURIUserInMemRepo) RemoveAllByUser(ctx context.Context, userID string) error {
	if userID == "" {
		return nil
	}

	toRemove, err := imsu.ListAllByUser(ctx, userID)
	if err != nil {
		return err
	}

	for _, entity := range toRemove {
		err := imsu.Remove(ctx, entity.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (imsu *ShortURIUserInMemRepo) RemoveAllByShortURI(ctx context.Context, shortURIID string) error {
	if shortURIID == "" {
		return nil
	}

	toRemove, err := imsu.ListAllByShortURI(ctx, shortURIID)
	if err != nil {
		return err
	}

	for _, entity := range toRemove {
		err := imsu.Remove(ctx, entity.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (imsu *ShortURIUserInMemRepo) Count(ctx context.Context) (int, error) {
	return len(imsu.cache.GetShortURIUserCache()), nil
}

func (imsu *ShortURIUserInMemRepo) UniqueCount(ctx context.Context) (int, error) {
	return countFunc(ctx, imsu, func(user *model.ShortURIUser) string {
		return user.UserID
	})
}

func countFunc[K comparable](ctx context.Context, imsu *ShortURIUserInMemRepo, selector func(user *model.ShortURIUser) K) (int, error) {
	// Делаем копию и быстренько освободлаем ресурс (мьютекс)
	// способ не очень по памяти, в качестве альтернативы батчами по n записей
	imsu.cache.GetShortURIUserRWMutex().RLock()
	cacheData := slices.Collect(maps.Values(imsu.cache.GetShortURIUserCache()))
	dataCopy := make([]*model.ShortURIUser, len(cacheData))
	copy(dataCopy, cacheData)
	imsu.cache.GetShortURIUserRWMutex().RUnlock()

	res := make(map[K]struct{})
	for _, entity := range dataCopy {
		if err := ctx.Err(); err != nil {
			return 0, err
		}

		res[selector(entity)] = struct{}{}
	}

	return len(res), nil
}

func (imsu *ShortURIUserInMemRepo) CreateStmt(ctx context.Context, stmt *sql.Stmt, entity *model.ShortURIUser) (*model.ShortURIUser, error) {
	return imsu.Create(ctx, entity)
}

func (imsu *ShortURIUserInMemRepo) ChangeStmt(ctx context.Context, stmt *sql.Stmt, entity *model.ShortURIUser) (*model.ShortURIUser, error) {
	return imsu.Change(ctx, entity)
}

func (imsu *ShortURIUserInMemRepo) DeleteStmt(ctx context.Context, stmt *sql.Stmt, ID string) error {
	return imsu.Delete(ctx, ID)
}

func (imsu *ShortURIUserInMemRepo) DeleteByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIID string) error {
	return imsu.DeleteByUnique(ctx, userID, shortURIID)
}

func (imsu *ShortURIUserInMemRepo) DeleteAllByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIIds []string) error {
	return imsu.DeleteAllByUnique(ctx, userID, shortURIIds)
}

func (imsu *ShortURIUserInMemRepo) DeleteAllByUserStmt(ctx context.Context, stmt *sql.Stmt, userID string) error {
	return imsu.DeleteAllByUser(ctx, userID)
}

func (imsu *ShortURIUserInMemRepo) DeleteAllByShortURIStmt(ctx context.Context, stmt *sql.Stmt, shortURIID string) error {
	return imsu.DeleteAllByShortURI(ctx, shortURIID)
}

func (imsu *ShortURIUserInMemRepo) RemoveStmt(ctx context.Context, stmt *sql.Stmt, ID string) error {
	return imsu.Remove(ctx, ID)
}

func (imsu *ShortURIUserInMemRepo) RemoveByUniqueStmt(ctx context.Context, stmt *sql.Stmt, userID string, shortURIID string) error {
	return imsu.RemoveByUnique(ctx, userID, shortURIID)
}

func (imsu *ShortURIUserInMemRepo) RemoveAllByUserStmt(ctx context.Context, stmt *sql.Stmt, userID string) error {
	return imsu.RemoveAllByUser(ctx, userID)
}

func (imsu *ShortURIUserInMemRepo) RemoveAllByShortURIStmt(ctx context.Context, stmt *sql.Stmt, shortURIID string) error {
	return imsu.RemoveAllByShortURI(ctx, shortURIID)
}

func (imsu *ShortURIUserInMemRepo) delete(ctx context.Context, entity *model.ShortURIUser) error {
	if entity == nil {
		return nil
	}

	entity.Deleted = true

	_, err := imsu.Change(ctx, entity)

	return err
}
