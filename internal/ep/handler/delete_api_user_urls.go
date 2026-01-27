package handler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

func (cr *AppChiRouter) deleteAPIUserUrls(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("deleteAPIUserUrls start")
	defer cr.log.Info("deleteAPIUserUrls finish")

	userInfo, err := auth.UserInfoFromRequestJWT(r)
	if err != nil {
		// 401
		if errors.As(err, &errs.AuthInfoAbsentErr) || errors.As(err, &errs.AuthInfoInvalidErr) {
			http.Error(rw, err.Error(), http.StatusUnauthorized)

			return
		}

		// 500
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	defer utils.CloseOnly(r.Body)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	go cr.batchDeleteAsync(userInfo, data)

	// 202
	rw.WriteHeader(http.StatusAccepted)
}

func (cr *AppChiRouter) batchDeleteAsync(userInfo *auth.UserInfo, data []byte) {
	err := cr.userFacade.BatchDelete(context.WithValue(context.Background(), auth.ContextUserInfo, userInfo), bytes.NewReader(data))
	if err != nil {
		cr.log.Errorf("error batchDeleteAsync with [%v]", err)
	}
}
