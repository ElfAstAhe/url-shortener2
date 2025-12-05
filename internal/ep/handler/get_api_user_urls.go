package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

func (cr *AppChiRouter) getAPIUserUrls(rw http.ResponseWriter, r *http.Request) {
	res, err := cr.userFacade.ListAllShortens(r.Context())
	if err != nil {
		// 401
		if errors.As(err, &errs.AuthUnauthorizedErr) || errors.As(err, &errs.AuthInfoAbsentErr) {
			http.Error(rw, err.Error(), http.StatusUnauthorized)

			return
		}

		// ToDo: NoContent implementation
		// ..
	}

	// 200
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	if err := enc.Encode(res); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
