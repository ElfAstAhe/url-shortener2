package auth

import (
	"context"
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type JWTAuthRetriever struct {
	log logger.Logger
}

func NewJWTAuthRetriever(log logger.Logger) *JWTAuthRetriever {
	return &JWTAuthRetriever{
		log: log.GetLogger("JWTAuthRetriever"),
	}
}

func (ar *JWTAuthRetriever) AuthRetriever(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ar.log.Info("AuthRetriever begin")
		defer ar.log.Info("AuthRetriever end")

		userInfo, err := auth.UserInfoFromRequestJWT(r)
		req := r
		if err != nil {
			ar.log.Warnf("err retrieve user info with error [%v]", err)
		} else {
			req = r.WithContext(context.WithValue(req.Context(), auth.ContextUserInfo, userInfo))
		}

		next.ServeHTTP(rw, req)
	})
}
