package audit

import (
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/audit"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/middleware"
	"github.com/ElfAstAhe/url-shortener2/pkg/client/audit/dto"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
	"github.com/ElfAstAhe/url-shortener2/pkg/utils"
)

type IncomeAuditMiddleware struct {
	watchPaths *middleware.PathMatchers
	publisher  audit.IncomePublisher
	log        logger.Logger
}

func NewIncomeAuditMiddleware(watchPaths *middleware.PathMatchers, publisher audit.IncomePublisher, log logger.Logger) *IncomeAuditMiddleware {
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
			pathMatcher := ia.watchPaths.GetPathMatcher(r.Method, r.RequestURI)
			if pathMatcher != nil {
				ia.log.Infof("IncomeAuditMiddleware.Handler matching path: method [%s] path [%s] pattern [%s]", pathMatcher.Method, pathMatcher.Path, pathMatcher.Pattern)
			}

			// custom response writer
			crw := middleware.NewCommonResponseWriter(rw)

			// serve request
			next.ServeHTTP(crw, r)

			// business logic
			if utils.IsSuccess(crw.Info.StatusCode) || utils.IsRedirection(crw.Info.StatusCode) {
				auditDto := ia.toDto(r, crw.Info)

				ia.log.Infof("Success, audit dto [%v], statusCode [%v]", auditDto, crw.Info.StatusCode)

				ia.publisher.NotifyAsync(auditDto)
			}
		} else {
			next.ServeHTTP(rw, r)
		}
	})
}

func (ia *IncomeAuditMiddleware) toDto(r *http.Request, rw *middleware.CommonResponseInfo) *dto.IncomeAuditDto {
	// ToDo: implement
	return &dto.IncomeAuditDto{}
}
