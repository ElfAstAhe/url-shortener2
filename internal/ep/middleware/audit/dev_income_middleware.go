package audit

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type DevHttpRequestDto struct {
	Host       string   `json:"host"`
	Method     string   `json:"method"`
	RequestURI string   `json:"requestURI"`
	Pattern    string   `json:"pattern"`
	URL        *url.URL `json:"url"`
	Proto      string   `json:"proto"`
	RemoteAddr string   `json:"remoteAddr"`
	Referer    string   `json:"referer"`
	UserAgent  string   `json:"userAgent"`
}

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
			req, err := json.Marshal(dim.requestToDto(r))
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

func (dim *DevIncomeMiddleware) requestToDto(r *http.Request) *DevHttpRequestDto {
	return &DevHttpRequestDto{
		Host:       r.Host,
		Method:     r.Method,
		RequestURI: r.RequestURI,
		Pattern:    r.URL.Path,
		URL:        r.URL,
		Proto:      r.Proto,
		RemoteAddr: r.RemoteAddr,
		Referer:    r.Referer(),
		UserAgent:  r.UserAgent(),
	}
}
