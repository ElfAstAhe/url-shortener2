package facade

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/dto"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/mapper"
	apperrs "github.com/ElfAstAhe/url-shortener2/internal/error"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

type ShortenFacadeImpl struct {
	service service.Shorter
	baseURL string
}

func NewShortenFacadeImpl(service service.Shorter, baseURL string) *ShortenFacadeImpl {
	return &ShortenFacadeImpl{
		service: service,
		baseURL: baseURL,
	}
}

func (sf *ShortenFacadeImpl) GetURL(ctx context.Context, req *http.Request) (string, error) {
	key := chi.URLParam(req, "key")
	// take raw data without user info
	data, err := sf.service.GetURL(ctx, key)
	if err != nil {
		return "", err
	}
	if data == "" {
		return "", nil
	}

	// get user info
	userInfo, errInfo := auth.UserInfoFromContext(ctx)
	if errInfo != nil {
		return data, nil
	}

	// get user data
	_, errUser := sf.service.GetURLUser(ctx, userInfo.UserID, key)
	if errUser != nil && !errors.As(errUser, &apperrs.DalSoftRemovedErr) {
		return "", errUser
	}

	return data, errUser
}

func (sf *ShortenFacadeImpl) CreateURL(ctx context.Context, income io.Reader) (*dto.ShortenCreateResponse, error) {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	cr, err := sf.getCRFromIncome(income)
	if err != nil {
		return nil, err
	}

	err = sf.validateCR(cr)
	if err != nil {
		return nil, err
	}

	key, err := sf.service.Store(ctx, userInfo.UserID, cr.URL)
	res, _ := mapper.ShortenCreateResponseFromKey(sf.baseURL, key)
	if err != nil {
		if res != nil {
			return res, err
		}

		return nil, err
	}

	return res, nil
}

func (sf *ShortenFacadeImpl) BatchCreateURL(ctx context.Context, income io.Reader) ([]*dto.ShortenBatchResponseItem, error) {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	incomeData, err := sf.getBatchFromIncome(income)
	if err != nil {
		return nil, err
	}

	serviceData, err := mapper.ShortenBatchFromDto(incomeData)
	if err != nil {
		return nil, err
	}

	serviceRes, err := sf.service.BatchStore(ctx, userInfo.UserID, serviceData)
	if err != nil {
		return nil, err
	}

	outcomeRes, err := mapper.ShortenBatchResponseFromKeys(serviceRes)
	if err != nil {
		return nil, err
	}

	return outcomeRes, nil
}

func (sf *ShortenFacadeImpl) Store(ctx context.Context, income io.Reader) (string, error) {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(income)
	if err != nil {
		return "", err
	}

	key, err := sf.service.Store(ctx, userInfo.UserID, string(data))
	if err != nil {
		return key, err
	}

	res, err := mapper.ShortenCreateResponseFromKey(sf.baseURL, key)
	if err != nil {
		return "", err
	}

	return res.Result, nil
}

func (sf *ShortenFacadeImpl) getCRFromIncome(income io.Reader) (*dto.ShortenCreateRequest, error) {
	dec := json.NewDecoder(income)
	var res = new(dto.ShortenCreateRequest)
	err := dec.Decode(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sf *ShortenFacadeImpl) validateCR(cr *dto.ShortenCreateRequest) error {
	return sf.validate("cr", cr.URL)
}

func (sf *ShortenFacadeImpl) getBatchFromIncome(income io.Reader) ([]*dto.ShortenBatchCreateItem, error) {
	dec := json.NewDecoder(income)
	var data = make([]*dto.ShortenBatchCreateItem, 0)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func (sf *ShortenFacadeImpl) validate(param string, data string) error {
	if strings.TrimSpace(data) == "" {
		return errs.NewAppInvalidArgumentError(param, "url is required")
	}

	return nil
}
