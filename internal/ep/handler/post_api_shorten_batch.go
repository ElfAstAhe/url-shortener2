package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/ElfAstAhe/url-shortener2/internal/utils"
)

func (cr *AppChiRouter) postAPIShortenBatch(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("postAPIShortenBatch start")
	defer cr.log.Info("postAPIShortenBatch finish")

	defer utils.CloseOnly(r.Body)

	res, err := cr.shortenFacade.BatchCreateURL(r.Context(), r.Body)
	if err != nil {
		cr.log.Error("postAPIShortenBatch err", zap.Error(err))
		// 500
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	// 201
	rw.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(rw)
	if err = enc.Encode(res); err != nil {
		cr.log.Error("postAPIShortenBatch err", zap.Error(err))
		// 500
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
