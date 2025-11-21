package service

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/repository"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	apperrs "github.com/ElfAstAhe/url-shortener2/internal/error"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

type ShorterImpl struct {
	shortURIRepo repository.ShortURIRepository
}

func NewShorterService(shortURIRepo repository.ShortURIRepository) (*ShorterImpl, error) {
	return &ShorterImpl{
		shortURIRepo: shortURIRepo,
	}, nil
}

// ShorterService

func (s *ShorterImpl) GetURL(ctx context.Context, key string) (string, error) {
	noAuthData, errNoAuth := s.getURLNoAuth(ctx, key)
	_, errAuth := s.getURLAuth(ctx, key)
	if errAuth != nil && errors.As(errAuth, &apperrs.DalSoftRemovedErr) {
		return "", errAuth
	}
	if errNoAuth != nil {
		return "", errNoAuth
	}

	return noAuthData, nil
}

func (s *ShorterImpl) getURLNoAuth(ctx context.Context, key string) (string, error) {
	res, err := s.shortURIRepo.GetByKey(ctx, key)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", nil
	}

	return res.OriginalURL.URL.String(), nil
}

func (s *ShorterImpl) getURLAuth(ctx context.Context, key string) (string, error) {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return "", err
	}
	if userInfo == nil {
		return "", errs.NewAuthInfoAbsentError("getURLAuth internal service method", nil)
	}
	res, err := s.shortURIRepo.GetByKeyUser(ctx, userInfo.UserID, key)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", nil
	}

	return res.OriginalURL.URL.String(), nil
}

func (s *ShorterImpl) Store(ctx context.Context, url string) (string, error) {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return "", err
	}

	key := utils.EncodeURIStr(url)
	res, err := model.NewShortURI(url, key)
	if err != nil {
		return "", err
	}

	res, err = s.shortURIRepo.Create(ctx, userInfo.UserID, res)
	if err != nil && res == nil {
		return "", err
	} else if err != nil {
		return res.Key, err
	}

	return res.Key, nil
}

func (s *ShorterImpl) BatchStore(ctx context.Context, source CorrelationUrls) (CorrelationShorts, error) {
	if len(source) == 0 {
		return CorrelationShorts{}, nil
	}

	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	repoBatch, err := s.toBatchSource(source)
	if err != nil {
		return nil, err
	}

	batchRes, err := s.shortURIRepo.BatchCreate(ctx, userInfo.UserID, repoBatch)
	if err != nil {
		return nil, err
	}

	res, err := s.toBatchResult(batchRes)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ShorterImpl) GetAllUserShorts(ctx context.Context, userID string) (UserShorts, error) {
	entities, err := s.shortURIRepo.ListAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	models, err := s.toUserShorts(entities)
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (s *ShorterImpl) BatchDelete(ctx context.Context, data UserBatchDeletes) error {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return err
	}

	return s.shortURIRepo.BatchDeleteByKeys(ctx, userInfo.UserID, data)
}

// ================

func (s *ShorterImpl) toBatchSource(source CorrelationUrls) (map[string]*model.ShortURI, error) {
	batch := make(map[string]*model.ShortURI)
	for correlation, origURL := range source {
		item, err := model.NewShortURI(origURL, utils.EncodeURIStr(origURL))
		if err != nil {
			return nil, err
		}
		batch[correlation] = item
	}

	return batch, nil
}

func (s *ShorterImpl) toBatchResult(source map[string]*model.ShortURI) (CorrelationShorts, error) {
	batch := make(CorrelationShorts)
	for correlation, shortURL := range source {
		batch[correlation] = utils.BuildNewURI(config.AppConfig.BaseURL, shortURL.Key)
	}

	return batch, nil
}

func (s *ShorterImpl) toUserShorts(entities []*model.ShortURI) (UserShorts, error) {
	if len(entities) == 0 {
		return nil, nil
	}
	res := make(UserShorts)
	for _, entity := range entities {
		res[entity.OriginalURL.URL.String()] = utils.BuildNewURI(config.AppConfig.BaseURL, entity.Key)
	}

	return res, nil
}
