package facade

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/dto"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/mapper"
	apperrs "github.com/ElfAstAhe/url-shortener2/internal/error"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
	"github.com/go-chi/chi/v5"
)

type ShortenFacadeImpl struct {
	service service.Shorter
}

func NewShortenFacadeImpl(service service.Shorter) *ShortenFacadeImpl {
	return &ShortenFacadeImpl{
		service: service,
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

func (sf *ShortenFacadeImpl) CreateURL(ctx context.Context, req *http.Request) (*dto.ShortenCreateResponse, error) {
	userInfo, err := auth.UserInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	cr, err := sf.getCRFromRequest(req)
	if err != nil {
		return nil, err
	}

	err = sf.validateCR(cr)
	if err != nil {
		return nil, err
	}

	key, err := sf.service.Store(ctx, userInfo.UserID, cr.URL)
	res, _ := mapper.ShortenCreateResponseFromKey(key)
	if err != nil {
		if res != nil {
			return res, err
		}

		return nil, err
	}

	return res, nil
}

func (sf *ShortenFacadeImpl) getCRFromRequest(req *http.Request) (*dto.ShortenCreateRequest, error) {
	defer utils.CloseOnly(req.Body)
	dec := json.NewDecoder(req.Body)
	var res = new(dto.ShortenCreateRequest)
	err := dec.Decode(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sf *ShortenFacadeImpl) validateCR(cr *dto.ShortenCreateRequest) error {
	if cr.URL == "" {
		return errs.NewAppInvalidArgumentError("cr", "url is required")
	}

	return nil
}
