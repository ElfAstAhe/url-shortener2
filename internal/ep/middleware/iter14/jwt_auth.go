package iter14

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type JWTAuthIter14 struct {
	WatchPaths *AuthIter14PathMatchers
	log        logger.Logger
}

func NewJWTAuthIter14(watchPaths *AuthIter14PathMatchers, logger logger.Logger) *JWTAuthIter14 {
	return &JWTAuthIter14{
		WatchPaths: watchPaths,
		log:        logger.GetLogger("Iter14Auth"),
	}
}

func (ai14 *JWTAuthIter14) Iter14Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ai14.log.Info("Iter14Auth begin")
		defer ai14.log.Info("Iter14Auth end")

		ai14.log.Info(fmt.Sprintf("request data info: Method [%s] Host [%s] Path [%s] URL [%v] Pattern [%s]", r.Method, r.Host, r.RequestURI, r.URL, r.Pattern))

		if ai14.WatchPaths.Match(r.Method, r.RequestURI) {
			matcher := ai14.WatchPaths.GetPathMatcher(r.Method, r.RequestURI)

			if matcher != nil {
				ai14.log.Info(fmt.Sprintf("Found matcher: method [%s] path [%s] pattern [%s] for request path [%s] ]", matcher.PathMatcher.Method, matcher.PathMatcher.Path, matcher.PathMatcher.Pattern, r.RequestURI))

				userInfo, err := auth.UserInfoFromRequestJWT(r)

				if !ai14.validateUserInfo(userInfo, err) {
					// answer with new auth cookie or error without cookie
					if err := ai14.answerError(userInfo, err, matcher, rw); err != nil {
						ai14.log.Errorf("got error on set auth cookie [%v]", err.Error())

						http.Error(rw, err.Error(), http.StatusInternalServerError)
					}

					// stop pipeline for request
					return
				}
			} else {
				http.Error(rw, "path matcher not found", http.StatusInternalServerError)

				return
			}
		}

		next.ServeHTTP(rw, r)
	})
}

func (ai14 *JWTAuthIter14) validateUserInfo(userInfo *auth.UserInfo, err error) bool {
	return userInfo != nil && err == nil
}

func (ai14 *JWTAuthIter14) answerError(userInfo *auth.UserInfo, err error, matcher *AuthIter14PathMatcher, rw http.ResponseWriter) error {
	uInfo := auth.BuildRandomUserInfo()
	tokenString, tbErr := auth.NewJWTStringFromUserInfo(uInfo)
	if tbErr != nil {
		return tbErr
	}
	var statusCode int
	var errMsg = ""
	switch {
	case errors.As(err, &errs.AuthInfoAbsentErr) && matcher.InfoAbsentStatusCode > 0:
		{
			ai14.setJWTAuthCookie(tokenString, rw)
			statusCode = matcher.InfoAbsentStatusCode
			errMsg = err.Error()
			break
		}
	case errors.As(err, &errs.AuthInfoInvalidErr) && matcher.InfoInvalidStatusCode > 0:
		{
			ai14.setJWTAuthCookie(tokenString, rw)
			statusCode = matcher.InfoInvalidStatusCode
			errMsg = err.Error()
			break
		}
	case err != nil && userInfo != nil && matcher.ErrAndUserInfoStatusCode > 0:
		{
			statusCode = matcher.ErrAndUserInfoStatusCode
			break
		}
	default:
		{
			statusCode = matcher.DefaultStatusCode
			break
		}
	}

	http.Error(rw, errMsg, statusCode)

	return nil
}

func (ai14 *JWTAuthIter14) setJWTAuthCookie(jwtString string, rw http.ResponseWriter) {
	http.SetCookie(rw, &http.Cookie{
		Name:     auth.CookieName,
		Value:    jwtString,
		SameSite: http.SameSiteStrictMode,
	})
}
