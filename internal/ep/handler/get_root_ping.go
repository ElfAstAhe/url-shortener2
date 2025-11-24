package handler

import "net/http"

func (cr *AppChiRouter) getRootPing(rw http.ResponseWriter, r *http.Request) {
	err := cr.toolFacade.Ping(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write([]byte("pong"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
