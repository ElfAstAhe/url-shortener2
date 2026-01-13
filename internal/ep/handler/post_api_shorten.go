package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/internal/ep/dto"
	apperrs "github.com/ElfAstAhe/url-shortener2/internal/error"
	"github.com/ElfAstAhe/url-shortener2/internal/utils"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

func (cr *AppChiRouter) postAPIShorten(rw http.ResponseWriter, r *http.Request) {
	cr.log.Info("postAPIShorten start")
	defer cr.log.Info("postAPIShorten finish")

	defer utils.CloseOnly(r.Body)

	res, err := cr.shortenFacade.CreateURL(r.Context(), r.Body)
	if err != nil {
		// 400
		if errors.As(err, &errs.AppInvalidArgumentErr) {
			http.Error(rw, err.Error(), http.StatusBadRequest)

			return
		}

		// 409
		if errors.As(err, &errs.ModelAlreadyExistsErr) || errors.As(err, &apperrs.BllConflictErr) {
			cr.sendRespAPIShortenCreate(res, http.StatusConflict, rw)

			return
		}

		// 500
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// 201
	cr.sendRespAPIShortenCreate(res, http.StatusCreated, rw)
}

// sendRespAPIShortenCreate send response with data and status code
func (cr *AppChiRouter) sendRespAPIShortenCreate(data *dto.ShortenCreateResponse, statusCode int, rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)

	enc := json.NewEncoder(rw)
	if err := enc.Encode(data); err != nil {
		// 500
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
