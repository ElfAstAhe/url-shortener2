package facade

import (
	"context"
	"errors"
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	apperrs "github.com/ElfAstAhe/url-shortener2/internal/error"
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

func (s *ShortenFacadeImpl) GetURL(ctx context.Context, req *http.Request) (string, error) {
	key := chi.URLParam(req, "key")
	// take raw data without user info
	data, err := s.service.GetURL(ctx, key)
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
	_, errUser := s.service.GetURLUser(ctx, userInfo.UserID, key)
	if errUser != nil && !errors.As(errUser, &apperrs.DalSoftRemovedErr) {
		return "", errUser
	}

	return data, errUser
}
