package service

import (
	"context"
	"strings"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/repository"
	apperrs "github.com/ElfAstAhe/url-shortener2/internal/error"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

type ShorterImpl struct {
	baseURL      string
	shortURIRepo repository.ShortURIRepository
}

func NewShorterService(config *config.Config, shortURIRepo repository.ShortURIRepository) *ShorterImpl {
	return &ShorterImpl{
		baseURL:      config.BaseURL,
		shortURIRepo: shortURIRepo,
	}
}

// ShorterService

func (s *ShorterImpl) GetURL(ctx context.Context, key string) (string, error) {
	res, err := s.shortURIRepo.GetByKey(ctx, key)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", nil
	}

	return res.OriginalURL.URL.String(), nil
}

func (s *ShorterImpl) GetURLUser(ctx context.Context, userID string, key string) (string, error) {
	if strings.TrimSpace(userID) == "" {
		return "", errs.NewAppInvalidArgumentError("userID", "must not be empty")
	}

	res, err := s.shortURIRepo.GetByKeyUser(ctx, userID, key)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", nil
	}

	return res.OriginalURL.URL.String(), nil
}

func (s *ShorterImpl) Store(ctx context.Context, userID string, url string) (string, error) {
	if strings.TrimSpace(userID) == "" {
		return "", errs.NewAuthInfoAbsentError("user id absent", nil)
	}

	key := utils.EncodeURIStr(url)
	res, err := model.NewShortURI(url, key)
	if err != nil {
		return "", err
	}

	res, err = s.shortURIRepo.Create(ctx, userID, res)
	if err != nil && res == nil {
		return "", err
	} else if err != nil {
		return res.Key, apperrs.NewBllConflictErrorEx(res.Key, err)
	}

	return res.Key, nil
}

func (s *ShorterImpl) BatchStore(ctx context.Context, userID string, source model.CorrelationUrls) (model.CorrelationShorts, error) {
	if len(source) == 0 {
		return make(model.CorrelationShorts), nil
	}

	if strings.TrimSpace(userID) == "" {
		return nil, errs.NewAppInvalidArgumentError("userID", "must not be empty")
	}

	repoBatch, err := s.toBatchSource(source)
	if err != nil {
		return nil, err
	}

	batchRes, err := s.shortURIRepo.BatchCreate(ctx, userID, repoBatch)
	if err != nil {
		return nil, err
	}

	res, err := s.toBatchResult(batchRes)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ShorterImpl) GetAllUserShorts(ctx context.Context, userID string) (model.UserShorts, error) {
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

func (s *ShorterImpl) BatchDelete(ctx context.Context, userID string, data model.UserBatchDeletes) error {
	if strings.TrimSpace(userID) == "" {
		return errs.NewAppInvalidArgumentError("userID", "is required")
	}

	return s.shortURIRepo.BatchDeleteByKeys(ctx, userID, data)
}

// ================

func (s *ShorterImpl) toBatchSource(source model.CorrelationUrls) (map[string]*model.ShortURI, error) {
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

func (s *ShorterImpl) toBatchResult(source map[string]*model.ShortURI) (model.CorrelationShorts, error) {
	batch := make(model.CorrelationShorts)
	for correlation, shortURL := range source {
		batch[correlation] = utils.BuildNewURI(s.baseURL, shortURL.Key)
	}

	return batch, nil
}

func (s *ShorterImpl) toUserShorts(entities []*model.ShortURI) (model.UserShorts, error) {
	if len(entities) == 0 {
		return nil, nil
	}
	res := make(model.UserShorts)
	for _, entity := range entities {
		res[entity.OriginalURL.URL.String()] = utils.BuildNewURI(s.baseURL, entity.Key)
	}

	return res, nil
}
