package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config/db"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/repository"
	apperrs "github.com/ElfAstAhe/url-shortener2/internal/error"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

type ShortURIInMemRepo struct {
	cache    db.InMemoryCache
	userRepo repository.ShortURIUserRepository
}

func NewShortURIInMemRepo(appDB db.DB, userRepo repository.ShortURIUserRepository) (*ShortURIInMemRepo, error) {
	if cache, ok := appDB.(db.InMemoryCache); ok {
		return &ShortURIInMemRepo{
			cache:    cache,
			userRepo: userRepo,
		}, nil
	}

	return nil, errs.NewAppInvalidConfigError("db", appDB, "db param does not implement InMemoryCache", nil)
}

func (ims *ShortURIInMemRepo) Get(ctx context.Context, id string) (*model.ShortURI, error) {
	ims.cache.GetShortURIRWMutex().RLock()
	defer ims.cache.GetShortURIRWMutex().RUnlock()
	res := ims.cache.GetShortURICache()[id]

	return res, nil
}

func (ims *ShortURIInMemRepo) GetByKey(ctx context.Context, key string) (*model.ShortURI, error) {
	if key == "" {
		return nil, nil
	}

	ims.cache.GetShortURIRWMutex().RLock()
	defer ims.cache.GetShortURIRWMutex().RUnlock()
	for _, value := range ims.cache.GetShortURICache() {
		if value.Key == key {
			return value, nil
		}
	}

	return nil, nil
}

func (ims *ShortURIInMemRepo) GetByKeyUser(ctx context.Context, userID string, key string) (*model.ShortURI, error) {
	if userID == "" {
		return nil, nil
	}
	if key == "" {
		return nil, nil
	}

	entity, err := ims.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}

	userLink, err := ims.userRepo.GetByUnique(ctx, userID, entity.ID)
	if err != nil {
		return nil, err
	}
	if userLink == nil {
		return nil, nil
	}
	if userLink.Deleted {
		return nil, apperrs.NewDalSoftRemovedError("short_uri", nil)
	}

	return entity, nil
}

func (ims *ShortURIInMemRepo) Create(ctx context.Context, userID string, entity *model.ShortURI) (*model.ShortURI, error) {
	if err := model.ValidateShortURI(entity); err != nil {
		return nil, err
	}

	find, err := ims.GetByKey(ctx, entity.Key)
	if err != nil {
		return nil, err
	}
	if find != nil {
		if err := ims.addUser(ctx, find.ID, userID); err != nil {
			return nil, err
		}

		return find, nil
	}

	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	entity.ID = newID.String()

	ims.cache.GetShortURIRWMutex().Lock()
	defer ims.cache.GetShortURIRWMutex().Unlock()
	ims.cache.GetShortURICache()[entity.ID] = entity

	if err := ims.addUser(ctx, entity.ID, userID); err != nil {
		return nil, err
	}

	return entity, nil
}

func (ims *ShortURIInMemRepo) BatchCreate(ctx context.Context, userID string, batch map[string]*model.ShortURI) (map[string]*model.ShortURI, error) {
	res := make(map[string]*model.ShortURI)
	if len(batch) == 0 {
		return res, nil
	}

	for correlation, item := range batch {
		data, err := ims.Create(ctx, userID, item)
		if err != nil && data == nil {
			return nil, err
		}
		res[correlation] = data
	}

	return res, nil
}

func (ims *ShortURIInMemRepo) ListAllByUser(ctx context.Context, userID string) ([]*model.ShortURI, error) {
	if userID == "" {
		return nil, nil
	}
	// all shorten ids by user
	entityUserLinks, err := ims.userRepo.ListAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return ims.listAllByLinks(ctx, entityUserLinks)
}

func (ims *ShortURIInMemRepo) ListAllByKeys(ctx context.Context, keys []string) ([]*model.ShortURI, error) {
	res := make([]*model.ShortURI, 0)
	if len(keys) == 0 {
		return res, nil
	}

	ims.cache.GetShortURIRWMutex().RLock()
	defer ims.cache.GetShortURIRWMutex().RUnlock()
	for key, value := range ims.cache.GetShortURICache() {
		if value.Key == key {
			res = append(res, value)
		}
	}

	return res, nil
}

func (ims *ShortURIInMemRepo) Delete(ctx context.Context, ID string, userID string) error {
	if userID == "" {
		return nil
	}
	if ID == "" {
		return nil
	}

	return ims.userRepo.DeleteByUnique(ctx, userID, ID)
}

func (ims *ShortURIInMemRepo) BatchDeleteByKeys(ctx context.Context, userID string, keys []string) error {
	if userID == "" || len(keys) == 0 {
		return nil
	}
	ids, err := ims.listIdsByKeys(ctx, keys)
	if err != nil {
		return err
	}

	inCh := ims.iter15Generator(ctx, ids)

	channels := ims.iter15FanOut(ctx, userID, inCh)

	res := ims.iter15FanIn(ctx, channels...)

	var eg errgroup.Group
	for dml := range res {
		eg.Go(func() error {
			return dml.Err
		})

	}

	return eg.Wait()
}

func (ims *ShortURIInMemRepo) iter15Generator(ctx context.Context, ids []string) chan string {
	inCh := make(chan string)

	go func() {
		defer close(inCh)

		for _, id := range ids {
			select {
			case <-ctx.Done():
				return
			case inCh <- id:
			}
		}
	}()

	return inCh
}

func (ims *ShortURIInMemRepo) iter15FanOut(ctx context.Context, userID string, inCh chan string) []chan *DMLResult {
	workersCount := 4
	res := make([]chan *DMLResult, workersCount)

	for index := 0; index < workersCount; index++ {
		res[index] = ims.iter15Delete(ctx, userID, inCh)
	}

	return res
}

func (ims *ShortURIInMemRepo) iter15Delete(ctx context.Context, userID string, inCh chan string) chan *DMLResult {
	res := make(chan *DMLResult)

	go func() {
		defer close(res)

		for id := range inCh {
			select {
			case <-ctx.Done():
				return
			case res <- NewDMLResult(ims.userRepo.DeleteByUnique(ctx, userID, id), "short_uri_users", id):
			}
		}
	}()

	return res
}

func (ims *ShortURIInMemRepo) iter15FanIn(ctx context.Context, resCh ...chan *DMLResult) chan *DMLResult {
	res := make(chan *DMLResult)

	var wg sync.WaitGroup

	for _, ch := range resCh {
		chClosure := ch
		wg.Add(1)

		go func() {
			defer wg.Done()

			for data := range chClosure {
				select {
				case <-ctx.Done():
					return
				case res <- data:
				}
			}
		}()
	}

	go func() {
		wg.Wait()

		close(res)
	}()

	return res
}

func (ims *ShortURIInMemRepo) listIdsByKeys(ctx context.Context, keys []string) ([]string, error) {
	res := make([]string, 0)
	if len(keys) == 0 {
		return res, nil
	}

	for _, key := range keys {
		entity, err := ims.GetByKey(ctx, key)
		if err != nil {
			return nil, err
		}
		if entity == nil {
			continue
		}
		res = append(res, entity.ID)
	}

	return res, nil
}

func (ims *ShortURIInMemRepo) listAllByLinks(ctx context.Context, userLinks []*model.ShortURIUser) ([]*model.ShortURI, error) {
	res := make([]*model.ShortURI, 0)
	if len(userLinks) == 0 {
		return res, nil
	}
	for _, userLink := range userLinks {
		find, err := ims.Get(ctx, userLink.ShortURIID)
		if err != nil {
			return nil, err
		}
		if find != nil {
			res = append(res, find)
		}
	}

	return res, nil
}

func (ims *ShortURIInMemRepo) addUser(ctx context.Context, ID string, userID string) error {
	find, err := ims.userRepo.GetByUnique(ctx, userID, ID)
	if err != nil {
		return err
	}
	if find != nil {
		return errs.NewModelAlreadyExistsError("short_uri_users", ID)
	}

	entity, err := model.NewShortURIUser(ID, userID)
	if err != nil {
		return err
	}

	res, err := ims.userRepo.Create(ctx, entity)

	if err != nil && res == nil {
		return err
	}

	return nil
}
