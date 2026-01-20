package handler

import (
	"encoding/json"
	"net/http"
)

func (cr *AppChiRouter) getAPIUserUrls(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getAPIUserUrls start")
	defer cr.log.Info("getAPIUserUrls finish")

	res, err := cr.userFacade.ListAllShortens(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// 200
	rw.WriteHeader(http.StatusOK)
	if len(res) == 0 {
		rw.WriteHeader(http.StatusNoContent)
	}
	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	if err := enc.Encode(res); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
