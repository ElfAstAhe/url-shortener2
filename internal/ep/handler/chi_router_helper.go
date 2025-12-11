package handler

import (
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
