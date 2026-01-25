package handler

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	apperrs "github.com/ElfAstAhe/url-shortener2/internal/error"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

func (cr *AppChiRouter) postRoot(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("postRoot start")
	defer cr.log.Info("postRoot finish")

	// Attention! Iter 14 only, specific auth business logic
	req, err := cr.setupIter14JWTOrSendError(rw, r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}
	defer utils.CloseOnly(req.Body)

	res, err := cr.shortenFacade.Store(req.Context(), req.Body)
	if err != nil && !errors.As(err, &apperrs.BllConflictErr) {
		// 500
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// outcome
	rw.Header().Set("Content-Type", "text/plain")
	if errors.As(err, &apperrs.BllConflictErr) {
		// 409
		rw.WriteHeader(http.StatusConflict)
	} else {
		// 201
		rw.WriteHeader(http.StatusCreated)
	}

	_, err = rw.Write([]byte(res))
	if err != nil {
		cr.log.Error("postAPIShortenBatch response write error", zap.Error(err))
	}
}

func (cr *AppChiRouter) setupIter14JWTOrSendError(rw http.ResponseWriter, r *http.Request) (*http.Request, error) {
	_, err := auth.UserInfoFromRequestJWT(r)
	if err != nil {
		if errors.As(err, &errs.AuthInfoInvalidErr) || errors.As(err, &errs.AuthInfoAbsentErr) {
			userInfo := auth.BuildRandomUserInfo()
			tokenString, err := auth.NewJWTStringFromUserInfo(userInfo)
			if err != nil {
				return nil, err
			}
			cr.setJWTAuthCookie(tokenString, rw)

			return r.WithContext(context.WithValue(r.Context(), auth.ContextUserInfo, userInfo)), nil
		}

		return nil, err
	}

	return r, nil
}
