package audit

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/audit"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	"github.com/ElfAstAhe/url-shortener2/pkg/client/audit/dto"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
	"github.com/ElfAstAhe/url-shortener2/pkg/utils"
)

type IncomeAuditMiddleware struct {
	watchPaths *IncomeAuditPathMatchers
	publisher  audit.IncomePublisher
	log        logger.Logger
}

func NewIncomeAuditMiddleware(watchPaths *IncomeAuditPathMatchers, publisher audit.IncomePublisher, log logger.Logger) *IncomeAuditMiddleware {
	return &IncomeAuditMiddleware{
		watchPaths: watchPaths,
		publisher:  publisher,
		log:        log.GetLogger("IncomeAuditMiddleware"),
	}
}

func (ia *IncomeAuditMiddleware) Audit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ia.log.Info("IncomeAuditMiddleware.Handler begin")
		defer ia.log.Info("IncomeAuditMiddleware.Handler end")

		if ia.watchPaths.Match(r.Method, r.RequestURI) {
			matcher := ia.watchPaths.GetPathMatcher(r.Method, r.RequestURI)
			if matcher != nil {
				ia.log.Infof("IncomeAuditMiddleware.Handler matching path: method [%s] path [%s] pattern [%s]", matcher.PathMatcher.Method, matcher.PathMatcher.Path, matcher.PathMatcher.Pattern)

				originalBody, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)

					return
				}

				r.Body = io.NopCloser(bytes.NewBuffer(originalBody))

				// custom response writer
				iarw := NewIncomeAuditResponseWriter(rw)

				// serve request
				next.ServeHTTP(iarw, r)

				// business logic
				if utils.IsSuccess(iarw.Info.StatusCode) || utils.IsRedirection(iarw.Info.StatusCode) {
					var data = make([]byte, 0)
					if utils.IsSuccess(iarw.Info.StatusCode) {
						data = append(data, originalBody...)
					} else if utils.IsRedirection(iarw.Info.StatusCode) {
						data = append(data, []byte(iarw.Header().Get("Location"))...)
					}

					if len(data) == 0 {
						ia.log.Warnf("income audit data is empty for: method [%s] path [%s] pattern [%s] request path [%s]", matcher.PathMatcher.Method, matcher.PathMatcher.Path, matcher.PathMatcher.Pattern, r.RequestURI)

						return
					}

					userInfo, err := auth.UserInfoFromContext(r.Context())
					userID := ""
					if err != nil {
						ia.log.Warnf("error extracting user info from request context: [%v]", zap.Error(err))
					} else {
						userID = userInfo.UserID
					}

					auditDto, err := ia.toDto(matcher, userID, data)
					if err != nil {
						ia.log.Warnf("error converting audit data to dto: [%v]", zap.Error(err))
					}

					ia.log.Infof("Success, audit dto [%v], statusCode [%v]", auditDto, iarw.Info.StatusCode)

					ia.publisher.NotifyAsync(auditDto)
				}
			} else {
				http.Error(rw, "path matcher not found", http.StatusInternalServerError)
			}
		} else {
			next.ServeHTTP(rw, r)
		}
	})
}

func (ia *IncomeAuditMiddleware) toDto(matcher *IncomeAuditPathMatcher, userID string, data []byte) (*dto.IncomeAuditDto, error) {
	url, err := matcher.urlExtractorFunc(data)
	if err != nil {
		return nil, err
	}

	return dto.NewIncomeAuditDto(time.Now(), ia.methodToAction(matcher.PathMatcher.Method), userID, url), nil
}

func (ia *IncomeAuditMiddleware) methodToAction(method string) string {
	switch method {
	case http.MethodGet:
		return dto.IncomeAuditActionFollow
	case http.MethodPost:
		return dto.IncomeAuditActionShorten
	}

	return fmt.Sprintf("unacceptable method [%s]", method)
}
