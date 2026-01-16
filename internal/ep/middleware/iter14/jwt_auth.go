package iter14

import (
	"fmt"
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type JWTAuthIter14 struct {
	WatchPaths *JWTAuthIter14Paths `json:"watch_paths"`
	log        logger.Logger
}

func NewJWTAuthIter14(watchPaths *JWTAuthIter14Paths, logger logger.Logger) *JWTAuthIter14 {
	return &JWTAuthIter14{
		WatchPaths: watchPaths,
		log:        logger.GetLogger("Iter14Auth"),
	}
}

func (a *JWTAuthIter14) Iter14Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		a.log.Info("Iter14Auth begin")
		defer a.log.Info("Iter14Auth end")

		a.log.Info(fmt.Sprintf("request data info: Method [%s] Host [%s] Path [%s] URL [%v] Pattern [%s]", r.Method, r.Host, r.RequestURI, r.URL, r.Pattern))

		// ToDo: implement

		next.ServeHTTP(rw, r)
	})
}
