package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

func (cr *AppChiRouter) getAPIInternalStats(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getInternalStats start")
	defer cr.log.Info("getInternalStats finish")

	res, err := cr.toolFacade.InternalStats(r.Context(), &r.Header)
	if err != nil {
		// 403
		if errors.As(err, &errs.AuthUnauthorizedErr) {
			http.Error(rw, err.Error(), http.StatusForbidden)

			return
		}

		// 500
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// 200
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(rw)
	if err := enc.Encode(res); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
