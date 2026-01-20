package handler

import (
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
)

func (cr *AppChiRouter) deleteAPIUserUrls(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("deleteAPIUserUrls start")
	defer cr.log.Info("deleteAPIUserUrls finish")

	if !auth.HasUserInfoInRequest(r) {
		// 401
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)

		return
	}

	go cr.batchDeleteAsync(r)

	// 202
	rw.WriteHeader(http.StatusAccepted)
}

func (cr *AppChiRouter) batchDeleteAsync(r *http.Request) {
	defer utils.CloseOnly(r.Body)

	err := cr.userFacade.BatchDelete(r.Context(), r.Body)
	if err != nil {
		cr.log.Errorf("error batchDeleteAsync with [%v]", err)
	}
}
