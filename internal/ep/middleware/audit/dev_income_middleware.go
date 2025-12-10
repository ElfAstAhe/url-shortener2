package audit

import (
	"encoding/json"
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type DevIncomeMiddleware struct {
	log     logger.Logger
	logging bool
}

func NewDevIncomeMiddleware(logger logger.Logger, logging bool) *DevIncomeMiddleware {
	return &DevIncomeMiddleware{
		log:     logger.GetLogger("dev audit"),
		logging: logging,
	}
}

func (dim *DevIncomeMiddleware) Handle(next http.Handler) http.Handler {
	meth := func(rw http.ResponseWriter, r *http.Request) {
		if dim.logging {
			dim.log.Infof("%s %s", r.Method, r.RequestURI)
			req, err := json.Marshal(r)
			if err != nil {
				dim.log.Errorf("dev http request audit error [%v]", err)
			} else {
				dim.log.Infof("dev http request audit [%s]", string(req))
			}
		}

		next.ServeHTTP(rw, r)
	}

	return http.HandlerFunc(meth)
}
