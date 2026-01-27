package mapper

import (
	"strings"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/model"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/dto"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
)

func ShortenCreateResponseFromKey(baseURL string, key string) (*dto.ShortenCreateResponse, error) {
	if strings.TrimSpace(key) == "" || strings.TrimSpace(baseURL) == "" {
		return nil, nil
	}

	return &dto.ShortenCreateResponse{
		Result: utils.BuildNewURI(baseURL, key),
	}, nil
}

func ShortenCreateResponseFromEntity(baseURL string, entity *model.ShortURI) (*dto.ShortenCreateResponse, error) {
	if entity == nil {
		return nil, nil
	}

	return ShortenCreateResponseFromKey(baseURL, entity.Key)
}

func ShortenBatchResponseFromKeys(source map[string]string) ([]*dto.ShortenBatchResponseItem, error) {
	if len(source) == 0 {
		return make([]*dto.ShortenBatchResponseItem, 0), nil
	}
	res := make([]*dto.ShortenBatchResponseItem, 0, len(source))
	for key, value := range source {
		res = append(res, &dto.ShortenBatchResponseItem{
			CorrelationID: key,
			ShortURL:      value,
		})
	}

	return res, nil
}

func ShortenBatchResponseFromEntity(baseURL string, source map[string]*model.ShortURI) ([]*dto.ShortenBatchResponseItem, error) {
	if len(source) == 0 {
		return make([]*dto.ShortenBatchResponseItem, 0), nil
	}
	res := make([]*dto.ShortenBatchResponseItem, 0, len(source))
	for key, value := range source {
		res = append(res, &dto.ShortenBatchResponseItem{
			CorrelationID: key,
			ShortURL:      utils.BuildNewURI(baseURL, value.Key),
		})
	}

	return res, nil
}

func ShortenBatchFromDto(source []*dto.ShortenBatchCreateItem) (map[string]string, error) {
	res := make(map[string]string)
	if len(source) == 0 {
		return res, nil
	}

	for _, item := range source {
		res[item.CorrelationID] = item.OriginalURL
	}

	return res, nil
}

func UserShortensFromModel(source map[string]string) ([]*dto.UserShorten, error) {
	if len(source) == 0 {
		return nil, nil
	}

	res := make([]*dto.UserShorten, 0, len(source))
	for key, value := range source {
		res = append(res, dto.NewUserShorten(value, key))
	}

	return res, nil
}

func UserBatchDeletesFromDto(source dto.ShortenBatchDeleteRequest) (model.UserBatchDeletes, error) {
	res := make(model.UserBatchDeletes, 0)
	if len(source) == 0 {
		return res, nil
	}

	for _, item := range source {
		res = append(res, item)
	}

	return res, nil
}
