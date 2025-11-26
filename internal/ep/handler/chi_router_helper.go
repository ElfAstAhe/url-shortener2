package handler

import (
	"context"
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
)

func (cr *AppChiRouter) setJWTAuthCookie(jwtString string, rw http.ResponseWriter) {
	http.SetCookie(rw, &http.Cookie{
		Name:     auth.CookieName,
		Value:    jwtString,
		SameSite: http.SameSiteStrictMode,
	})
}

func (cr *AppChiRouter) hasUserInfo(ctx context.Context) bool {
	userInfo, err := auth.UserInfoFromContext(ctx)

	return err == nil && userInfo != nil
}
