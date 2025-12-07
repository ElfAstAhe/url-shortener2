package handler

import (
	"encoding/json"
	"net/http"
)

func (cr *AppChiRouter) postAPIShortenBatch(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("postAPIShortenBatch start")
	defer cr.log.Info("postAPIShortenBatch finish")

	res, err := cr.shortenFacade.BatchCreateURL(r.Context(), r)
	if err != nil {
		// 500
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(rw)
	if err = enc.Encode(res); err != nil {
		// 500
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
