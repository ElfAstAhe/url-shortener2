package handler

import (
	"errors"
	"net/http"

	apperrs "github.com/ElfAstAhe/url-shortener2/internal/error"
)

func (cr *AppChiRouter) getRoot(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("getRoot start")
	defer cr.log.Info("getRoot finish")

	res, err := cr.shortenFacade.GetURL(r.Context(), r)
	if err != nil {
		// 410
		if errors.As(err, &apperrs.DalSoftRemovedErr) {
			http.Error(rw, err.Error(), http.StatusGone)

			return
		}

		// 500
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// 404
	if res == "" {
		http.Error(rw, "No shorter url found!", http.StatusNotFound)

		return
	}

	// 307
	rw.Header().Set("Content-Type", "plain/text")
	http.Redirect(rw, r, res, http.StatusTemporaryRedirect)
}
