package audit

import (
	"net/http"

	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/audit"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/middleware"
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
		crw := middleware.NewCommonResponseWriter(rw)

		next.ServeHTTP(crw, r)

		if utils.IsSuccess(crw.Info.StatusCode) || utils.IsRedirection(crw.Info.StatusCode) {
			// ToDo: implement
			// ..
		}
	})
}
