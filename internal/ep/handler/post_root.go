package handler

import (
	"errors"
	"net/http"

	apperrs "github.com/ElfAstAhe/url-shortener2/internal/error"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
	"go.uber.org/zap"
)

func (cr *AppChiRouter) postRoot(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("postAPIShortenBatch start")
	defer cr.log.Info("postAPIShortenBatch finish")

	defer utils.CloseOnly(r.Body)

	res, err := cr.shortenFacade.Store(r.Context(), r.Body)
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
